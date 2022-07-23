package warlock

import (
	"strconv"
	"time"

	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/stats"
)

func (warlock *Warlock) registerImmolateSpell() {
	actionID := core.ActionID{SpellID: 47811}
	baseCost := 0.17 * warlock.BaseMana
	costReductionFactor := 1.0
	if float64(warlock.Talents.Cataclysm) > 0 {
		costReductionFactor -= 0.01 + 0.03*float64(warlock.Talents.Cataclysm)
	}
	spellCoefficient:= 0.2 + 0.04*float64(warlock.Talents.ShadowAndFlame)

	effect := core.SpellEffect{
		BonusSpellCritRating: core.TernaryFloat64(warlock.Talents.Devastation, 1, 0) * 5 * core.CritRatingPerCritChance,
		DamageMultiplier:     1,
		ThreatMultiplier:     1 - 0.1*float64(warlock.Talents.DestructiveReach),
		BaseDamage:           core.BaseDamageConfigMagic(460.0, 460.0, spellCoefficient),
		OutcomeApplier:       warlock.OutcomeFuncMagicHitAndCrit(warlock.SpellCritMultiplier(1, float64(warlock.Talents.Ruin)/5)),
		OnSpellHitDealt:      applyDotOnLanded(&warlock.ImmolateDot),
		ProcMask:             core.ProcMaskSpellDamage,
	}

	warlock.Immolate = warlock.RegisterSpell(core.SpellConfig{
		ActionID:     actionID,
		SpellSchool:  core.SpellSchoolFire,
		ResourceType: stats.Mana,
		BaseCost:     baseCost,

		Cast: core.CastConfig{
			DefaultCast: core.Cast{
				Cost:     baseCost * costReductionFactor,
				GCD:      core.GCDDefault,
				CastTime: time.Millisecond * (2000 - 100*time.Duration(warlock.Talents.Bane)),
			},
			ModifyCast: func(_ *core.Simulation, _ *core.Spell, cast *core.Cast) {
				cast.GCD = time.Duration(float64(cast.GCD) * warlock.backdraftModifier())
				cast.CastTime = time.Duration(float64(cast.CastTime) * warlock.backdraftModifier())
			},
		},
		ApplyEffects: core.ApplyEffectFuncDirectDamage(effect),
	})

	target := warlock.CurrentTarget
	applier := warlock.OutcomeFuncTick()
	if warlock.Talents.Pandemic {
		applier = warlock.OutcomeFuncMagicCrit(warlock.SpellCritMultiplier(1, 1))
	}

	warlock.ImmolateDot = core.NewDot(core.Dot{
		Spell: warlock.Immolate,
		Aura: target.RegisterAura(core.Aura{
			Label:    "immolate-" + strconv.Itoa(int(warlock.Index)),
			ActionID: actionID,
		}),
		NumberOfTicks: 5 + int(warlock.Talents.MoltenCore)*3 +
			core.TernaryInt(warlock.HasSetBonus(ItemSetVoidheartRaiment, 4), 1, 0), // voidheart 4p gives 1 extra tick
		TickLength: time.Second * 3,
		TickEffects: core.TickFuncSnapshot(target, core.SpellEffect{
			DamageMultiplier: 1,
			ThreatMultiplier: 1 - 0.1*float64(warlock.Talents.DestructiveReach),
			BaseDamage:       core.BaseDamageConfigMagicNoRoll(785/5, spellCoefficient),
			OutcomeApplier:   applier,
			IsPeriodic:       true,
			ProcMask:         core.ProcMaskPeriodicDamage,
		}),
	})
}

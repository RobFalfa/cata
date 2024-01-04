package dps

import (
	"github.com/wowsims/wotlk/sim/core"
	"github.com/wowsims/wotlk/sim/core/proto"
	"github.com/wowsims/wotlk/sim/deathknight"
)

func RegisterDpsDeathknight() {
	core.RegisterAgentFactory(
		proto.Player_Deathknight{},
		proto.Spec_SpecDeathknight,
		func(character *core.Character, options *proto.Player) core.Agent {
			return NewDpsDeathknight(character, options)
		},
		func(player *proto.Player, spec interface{}) {
			playerSpec, ok := spec.(*proto.Player_Deathknight)
			if !ok {
				panic("Invalid spec value for Deathknight!")
			}
			player.Spec = playerSpec
		},
	)
}

type DpsDeathknight struct {
	*deathknight.Deathknight
}

func NewDpsDeathknight(character *core.Character, player *proto.Player) *DpsDeathknight {
	dk := player.GetDeathknight()

	dpsDk := &DpsDeathknight{
		Deathknight: deathknight.NewDeathknight(character, deathknight.DeathknightInputs{
			StartingRunicPower:  dk.Options.StartingRunicPower,
			PrecastGhoulFrenzy:  dk.Options.PrecastGhoulFrenzy,
			PrecastHornOfWinter: dk.Options.PrecastHornOfWinter,
			PetUptime:           dk.Options.PetUptime,
			DrwPestiApply:       dk.Options.DrwPestiApply,
			BloodOpener:         dk.Rotation.BloodOpener,
			IsDps:               true,
			NewDrw:              true,
			DiseaseDowntime:     dk.Options.DiseaseDowntime,
			VirulenceRefresh:    dk.Rotation.VirulenceRefresh,

			RefreshHornOfWinter: dk.Rotation.RefreshHornOfWinter,
			ArmyOfTheDeadType:   dk.Rotation.ArmyOfTheDead,
			StartingPresence:    dk.Rotation.StartingPresence,
			UseAMS:              dk.Rotation.UseAms,
			AvgAMSSuccessRate:   dk.Rotation.AvgAmsSuccessRate,
			AvgAMSHit:           dk.Rotation.AvgAmsHit,
		}, player.TalentsString, dk.Rotation.PreNerfedGargoyle),
	}

	dpsDk.Inputs.UnholyFrenzyTarget = dk.Options.UnholyFrenzyTarget

	dpsDk.EnableAutoAttacks(dpsDk, core.AutoAttackOptions{
		MainHand:       dpsDk.WeaponFromMainHand(dpsDk.DefaultMeleeCritMultiplier()),
		OffHand:        dpsDk.WeaponFromOffHand(dpsDk.DefaultMeleeCritMultiplier()),
		AutoSwingMelee: true,
	})

	return dpsDk
}

func (dk *DpsDeathknight) FrostPointsInBlood() int32 {
	return dk.Talents.Butchery + dk.Talents.Subversion + dk.Talents.BladeBarrier + dk.Talents.DarkConviction
}

func (dk *DpsDeathknight) FrostPointsInUnholy() int32 {
	return dk.Talents.ViciousStrikes + dk.Talents.Virulence + dk.Talents.Epidemic + dk.Talents.RavenousDead + dk.Talents.Necrosis + dk.Talents.BloodCakedBlade
}

func (dk *DpsDeathknight) GetDeathknight() *deathknight.Deathknight {
	return dk.Deathknight
}

func (dk *DpsDeathknight) Initialize() {
	dk.Deathknight.Initialize()
}

func (dk *DpsDeathknight) Reset(sim *core.Simulation) {
	dk.Deathknight.Reset(sim)
	dk.Presence = deathknight.UnsetPresence
}

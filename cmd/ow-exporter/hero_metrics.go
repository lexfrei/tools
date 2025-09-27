//nolint:govet,lll,tagalign // Multi-line struct tags and alignment are intentionally used for readability per CLAUDE.md guidelines
package main

import (
	"reflect"
	"time"
)

// Constants for frequently used strings.
const (
	CountMetricType                 = "count"
	DurationMetricType              = "duration"
	PercentageMetricType            = "percentage"
	NumberMetricType                = "number"
	MouseKeyboardViewActiveSelector = ".mouseKeyboard-view.is-active"
	ControllerViewActiveSelector    = ".controller-view.is-active"
	QuickPlayViewActiveSelector     = ".quickPlay-view.is-active"
	CompetitiveViewActiveSelector   = ".competitive-view.is-active"
)

// CommonMetrics defines the 15 core metrics available for all heroes.
type CommonMetrics struct {
	TimePlayed time.Duration `ow:"time_played"
		prometheus:"ow_hero_time_played_seconds"
		help:"Total time played on hero"
		path:"[data-category-id='0x0860000000000021']"
		type:"duration"`

	GamesWon int `ow:"games_won"
		prometheus:"ow_hero_games_won_total"
		help:"Total number of games won with hero"
		path:"[data-category-id='0x0860000000000039']"
		type:"number"`

	WinPercentage float64 `ow:"win_percentage"
		prometheus:"ow_hero_win_percentage"
		help:"Win percentage with hero"
		path:"[data-category-id='0x08600000000003D1']"
		type:"percentage"`

	WeaponAccuracy float64 `ow:"weapon_accuracy"
		prometheus:"ow_hero_weapon_accuracy_percent"
		help:"Best weapon accuracy percentage with hero"
		path:"[data-category-id='0x08600000000001BB']"
		type:"percentage"`

	EliminationsPerLife float64 `ow:"eliminations_per_life"
		prometheus:"ow_hero_eliminations_per_life"
		help:"Average eliminations per life with hero"
		path:"[data-category-id='0x08600000000003D2']"
		type:"number"`

	KillStreakBest int `ow:"kill_streak_best"
		prometheus:"ow_hero_kill_streak_best"
		help:"Best kill streak achieved with hero"
		path:"[data-category-id='0x0860000000000223']"
		type:"number"`

	MultikillBest int `ow:"multikill_best"
		prometheus:"ow_hero_multikill_best"
		help:"Best multikill achieved with hero"
		path:"[data-category-id='0x0860000000000346']"
		type:"number"`

	EliminationsPer10Min float64 `ow:"eliminations_per_10min"
		prometheus:"ow_hero_eliminations_per_10min_avg"
		help:"Average eliminations per 10 minutes with hero"
		path:"[data-category-id='0x08600000000004D4']"
		type:"number"`

	DeathsPer10Min float64 `ow:"deaths_per_10min"
		prometheus:"ow_hero_deaths_per_10min_avg"
		help:"Average deaths per 10 minutes with hero"
		path:"[data-category-id='0x08600000000004D3']"
		type:"number"`

	FinalBlowsPer10Min float64 `ow:"final_blows_per_10min"
		prometheus:"ow_hero_final_blows_per_10min_avg"
		help:"Average final blows per 10 minutes with hero"
		path:"[data-category-id='0x08600000000004D5']"
		type:"number"`

	SoloKillsPer10Min float64 `ow:"solo_kills_per_10min"
		prometheus:"ow_hero_solo_kills_per_10min_avg"
		help:"Average solo kills per 10 minutes with hero"
		path:"[data-category-id='0x08600000000004DA']"
		type:"number"`

	ObjectiveKillsPer10Min float64 `ow:"objective_kills_per_10min"
		prometheus:"ow_hero_objective_kills_per_10min_avg"
		help:"Average objective kills per 10 minutes with hero"
		path:"[data-category-id='0x08600000000004D8']"
		type:"number"`

	ObjectiveTimePer10Min time.Duration `ow:"objective_time_per_10min"
		prometheus:"ow_hero_objective_time_per_10min_avg"
		help:"Average objective time per 10 minutes with hero"
		path:"[data-category-id='0x08600000000004D9']"
		type:"duration"`

	HeroDamagePer10Min int64 `ow:"hero_damage_per_10min"
		prometheus:"ow_hero_damage_per_10min_avg"
		help:"Average hero damage per 10 minutes"
		path:"[data-category-id='0x08600000000004BD']"
		type:"number"`

	HealingPer10Min int64 `ow:"healing_per_10min"
		prometheus:"ow_hero_healing_per_10min_avg"
		help:"Average healing done per 10 minutes"
		path:"[data-category-id='0x08600000000004D6']"
		type:"number"`
}

// Soldier76Metrics defines Soldier: 76 specific metrics.
type Soldier76Metrics struct {
	CommonMetrics // Embedded common metrics

	HelixRocketKills int `ow:"helix_rocket_kills"
		prometheus:"ow_hero_helix_rocket_kills_total"
		help:"Total eliminations with helix rockets"
		path:"[data-stat='helix_rocket_kills']"
		type:"number"`

	HelixRocketKillsBest int `ow:"helix_rocket_kills_best"
		prometheus:"ow_hero_helix_rocket_kills_best"
		help:"Most helix rocket kills in a single game"
		path:"[data-stat='helix_rocket_kills_best']"
		type:"number"`

	BioticFieldHealing int64 `ow:"biotic_field_healing"
		prometheus:"ow_hero_biotic_field_healing_total"
		help:"Total healing provided by biotic field"
		path:"[data-stat='biotic_field_healing']"
		type:"number"`

	TacticalVisorKills int `ow:"tactical_visor_kills"
		prometheus:"ow_hero_tactical_visor_kills_total"
		help:"Total eliminations during tactical visor ultimate"
		path:"[data-stat='tactical_visor_kills']"
		type:"number"`
}

// WidowmakerMetrics defines Widowmaker specific metrics.
type WidowmakerMetrics struct {
	CommonMetrics // Embedded common metrics

	ScopedAccuracy float64 `ow:"scoped_accuracy"
		prometheus:"ow_hero_scoped_accuracy_percent"
		help:"Scoped weapon accuracy percentage"
		path:"[data-stat='scoped_accuracy']"
		type:"percentage"`

	ScopedCriticalHits int `ow:"scoped_critical_hits"
		prometheus:"ow_hero_scoped_critical_hits_total"
		help:"Total scoped critical hits"
		path:"[data-stat='scoped_critical_hits']"
		type:"number"`

	VenomMineKills int `ow:"venom_mine_kills"
		prometheus:"ow_hero_venom_mine_kills_total"
		help:"Total eliminations with venom mine"
		path:"[data-stat='venom_mine_kills']"
		type:"number"`

	InfraSightAssists int `ow:"infra_sight_assists"
		prometheus:"ow_hero_infra_sight_assists_total"
		help:"Team assists provided by infra-sight ultimate"
		path:"[data-stat='infra_sight_assists']"
		type:"number"`
}

// GenjiMetrics defines Genji specific metrics.
type GenjiMetrics struct {
	CommonMetrics // Embedded common metrics

	DeflectionKills  int   `ow:"deflection_kills" prometheus:"ow_hero_deflection_kills_total" help:"Total eliminations with deflected projectiles" path:"[data-stat='deflection_kills']" type:"number"`
	SwiftStrikeKills int   `ow:"swift_strike_kills" prometheus:"ow_hero_swift_strike_kills_total" help:"Total eliminations with swift strike" path:"[data-stat='swift_strike_kills']" type:"number"`
	DragonbladeKills int   `ow:"dragonblade_kills" prometheus:"ow_hero_dragonblade_kills_total" help:"Total eliminations during dragonblade ultimate" path:"[data-stat='dragonblade_kills']" type:"number"`
	DamageDeflected  int64 `ow:"damage_deflected" prometheus:"ow_hero_damage_deflected_total" help:"Total damage deflected by deflect ability" path:"[data-stat='damage_deflected']" type:"number"`
}

// TorbjornMetrics defines Torbjörn specific metrics.
type TorbjornMetrics struct {
	CommonMetrics // Embedded common metrics

	TurretKills       int   `ow:"turret_kills" prometheus:"ow_hero_turret_kills_total" help:"Total eliminations by deployed turret" path:"[data-stat='turret_kills']" type:"number"`
	TurretDamage      int64 `ow:"turret_damage" prometheus:"ow_hero_turret_damage_total" help:"Total damage dealt by turret" path:"[data-stat='turret_damage']" type:"number"`
	HammerKills       int   `ow:"hammer_kills" prometheus:"ow_hero_hammer_kills_total" help:"Total eliminations with forge hammer" path:"[data-stat='hammer_kills']" type:"number"`
	ArmorPacksCreated int   `ow:"armor_packs_created" prometheus:"ow_hero_armor_packs_created_total" help:"Total armor packs created for teammates" path:"[data-stat='armor_packs_created']" type:"number"`
	MoltenCoreKills   int   `ow:"molten_core_kills" prometheus:"ow_hero_molten_core_kills_total" help:"Eliminations during molten core ultimate" path:"[data-stat='molten_core_kills']" type:"number"`
}

// MercyMetrics defines Mercy specific metrics.
type MercyMetrics struct {
	CommonMetrics // Embedded common metrics

	DamageAmplified   int64 `ow:"damage_amplified" prometheus:"ow_hero_damage_amplified_total" help:"Total damage amplified with damage boost" path:"[data-stat='damage_amplified']" type:"number"`
	Resurrects        int   `ow:"resurrects" prometheus:"ow_hero_resurrects_total" help:"Total number of resurrections performed" path:"[data-stat='resurrects']" type:"number"`
	PlayersRezed      int   `ow:"players_rezed" prometheus:"ow_hero_players_rezed_total" help:"Total players resurrected" path:"[data-stat='players_rezed']" type:"number"`
	ValkyrieDamageAmp int64 `ow:"valkyrie_damage_amp" prometheus:"ow_hero_valkyrie_damage_amp_total" help:"Damage amplified during valkyrie ultimate" path:"[data-stat='valkyrie_damage_amp']" type:"number"`
	ValkyrieHealing   int64 `ow:"valkyrie_healing" prometheus:"ow_hero_valkyrie_healing_total" help:"Healing done during valkyrie ultimate" path:"[data-stat='valkyrie_healing']" type:"number"`
}

// ReinhardtMetrics defines Reinhardt specific metrics.
type ReinhardtMetrics struct {
	CommonMetrics // Embedded common metrics

	DamageBlocked     int64 `ow:"damage_blocked" prometheus:"ow_hero_damage_blocked_total" help:"Total damage blocked by barrier shield" path:"[data-stat='damage_blocked']" type:"number"`
	ChargeKills       int   `ow:"charge_kills" prometheus:"ow_hero_charge_kills_total" help:"Total eliminations with charge ability" path:"[data-stat='charge_kills']" type:"number"`
	FireStrikeKills   int   `ow:"fire_strike_kills" prometheus:"ow_hero_fire_strike_kills_total" help:"Total eliminations with fire strike" path:"[data-stat='fire_strike_kills']" type:"number"`
	EarthshatterKills int   `ow:"earthshatter_kills" prometheus:"ow_hero_earthshatter_kills_total" help:"Eliminations during earthshatter ultimate" path:"[data-stat='earthshatter_kills']" type:"number"`
	RocketHammerKills int   `ow:"rocket_hammer_kills" prometheus:"ow_hero_rocket_hammer_kills_total" help:"Total eliminations with rocket hammer" path:"[data-stat='rocket_hammer_kills']" type:"number"`
}

// AnaMetrics defines Ana specific metrics.
type AnaMetrics struct {
	CommonMetrics // Embedded common metrics

	ScopedAccuracy     float64 `ow:"scoped_accuracy" prometheus:"ow_hero_scoped_accuracy_percent" help:"Scoped weapon accuracy percentage" path:"[data-stat='scoped_accuracy']" type:"percentage"`
	UnscopedAccuracy   float64 `ow:"unscoped_accuracy" prometheus:"ow_hero_unscoped_accuracy_percent" help:"Unscoped weapon accuracy percentage" path:"[data-stat='unscoped_accuracy']" type:"percentage"`
	EnemiesSlept       int     `ow:"enemies_slept" prometheus:"ow_hero_enemies_slept_total" help:"Total enemies put to sleep with sleep dart" path:"[data-stat='enemies_slept']" type:"number"`
	BioticGrenadeKills int     `ow:"biotic_grenade_kills" prometheus:"ow_hero_biotic_grenade_kills_total" help:"Total eliminations with biotic grenade" path:"[data-stat='biotic_grenade_kills']" type:"number"`
	NanoboostAssists   int     `ow:"nanoboost_assists" prometheus:"ow_hero_nanoboost_assists_total" help:"Eliminations assisted by nanoboost ultimate" path:"[data-stat='nanoboost_assists']" type:"number"`
}

// DVaMetrics defines D.Va specific metrics.
type DVaMetrics struct {
	CommonMetrics // Embedded common metrics

	MechKills         int   `ow:"mech_kills" prometheus:"ow_hero_mech_kills_total" help:"Total eliminations while in mech" path:"[data-stat='mech_kills']" type:"number"`
	PilotKills        int   `ow:"pilot_kills" prometheus:"ow_hero_pilot_kills_total" help:"Total eliminations while out of mech" path:"[data-stat='pilot_kills']" type:"number"`
	DamageBlocked     int64 `ow:"damage_blocked" prometheus:"ow_hero_damage_blocked_total" help:"Total damage blocked by defense matrix" path:"[data-stat='damage_blocked']" type:"number"`
	SelfDestructKills int   `ow:"self_destruct_kills" prometheus:"ow_hero_self_destruct_kills_total" help:"Eliminations with self-destruct ultimate" path:"[data-stat='self_destruct_kills']" type:"number"`
	CallMechKills     int   `ow:"call_mech_kills" prometheus:"ow_hero_call_mech_kills_total" help:"Eliminations by calling down mech" path:"[data-stat='call_mech_kills']" type:"number"`
}

// GenerateMetricDefs uses reflection to generate MetricDef from a hero struct (defaults to PC QuickPlay).
func GenerateMetricDefs(heroStruct any) map[string]MetricDef {
	return GenerateMetricDefsWithContext(heroStruct, PlatformPC, GameModeQuickPlay)
}

// GenerateMetricDefsWithContext uses reflection to generate MetricDef with platform/gamemode context.
func GenerateMetricDefsWithContext(heroStruct any, platform Platform, gameMode GameMode) map[string]MetricDef {
	t := reflect.TypeOf(heroStruct)
	metrics := make(map[string]MetricDef)

	for i := range t.NumField() {
		field := t.Field(i)

		// Skip embedded CommonMetrics (will be handled separately)
		if field.Anonymous {
			continue
		}

		owTag := field.Tag.Get("ow")
		if owTag == "" {
			continue
		}

		// Generate platform/gamemode aware selector
		baseSelector := field.Tag.Get("path")
		selector := generatePlatformSelector(baseSelector, platform, gameMode)

		metrics[owTag] = MetricDef{
			PrometheusName: field.Tag.Get("prometheus"),
			Help:           field.Tag.Get("help"),
			Selector:       selector,
			ValueType:      field.Tag.Get("type"),
			Unit:           inferUnit(field.Tag.Get("type")),
		}
	}

	return metrics
}

// inferUnit determines the unit based on value type.
func inferUnit(valueType string) string {
	switch valueType {
	case DurationMetricType:
		return "seconds"
	case PercentageMetricType:
		return "percent"
	case NumberMetricType:
		return CountMetricType
	default:
		return CountMetricType
	}
}

// HeroMetricsRegistry maps hero IDs to their metric generation functions.
var HeroMetricsRegistry = map[string]func() any{
	// Initial 8 heroes
	"soldier-76": func() any { return Soldier76Metrics{} },
	"widowmaker": func() any { return WidowmakerMetrics{} },
	"genji":      func() any { return GenjiMetrics{} },
	"torbjorn":   func() any { return TorbjornMetrics{} },
	"mercy":      func() any { return MercyMetrics{} },
	"reinhardt":  func() any { return ReinhardtMetrics{} },
	"ana":        func() any { return AnaMetrics{} },
	"dva":        func() any { return DVaMetrics{} },

	// Support heroes
	"illari":     func() any { return IllariMetrics{} },
	"lifeweaver": func() any { return LifeweaverMetrics{} },
	"kiriko":     func() any { return KirikoMetrics{} },
	"baptiste":   func() any { return BaptisteMetrics{} },
	"lucio":      func() any { return LucioMetrics{} },
	"zenyatta":   func() any { return ZenyattaMetrics{} },
	"brigitte":   func() any { return BrigitteMetrics{} },

	// DPS heroes
	"cassidy": func() any { return CassidyMetrics{} },
	"tracer":  func() any { return TracerMetrics{} },
	"pharah":  func() any { return PharahMetrics{} },
	"sojourn": func() any { return SojournMetrics{} },
	"mei":     func() any { return MeiMetrics{} },
	"junkrat": func() any { return JunkratMetrics{} },
	"reaper":  func() any { return ReaperMetrics{} },
	"hanzo":   func() any { return HanzoMetrics{} },

	// Tank heroes
	"winston":       func() any { return WinstonMetrics{} },
	"roadhog":       func() any { return RoadhogMetrics{} },
	"zarya":         func() any { return ZaryaMetrics{} },
	"mauga":         func() any { return MaugaMetrics{} },
	"hazard":        func() any { return HazardMetrics{} },
	"junker-queen":  func() any { return JunkerQueenMetrics{} },
	"orisa":         func() any { return OrisaMetrics{} },
	"sigma":         func() any { return SigmaMetrics{} },
	"wrecking-ball": func() any { return WreckingBallMetrics{} },

	// Additional DPS heroes
	"doomfist": func() any { return DoomfistMetrics{} },
	"sombra":   func() any { return SombraMetrics{} },
	"symmetra": func() any { return SymmetraMetrics{} },
	"bastion":  func() any { return BastionMetrics{} },
	"ashe":     func() any { return AsheMetrics{} },
	"echo":     func() any { return EchoMetrics{} },
	"venture":  func() any { return VentureMetrics{} },

	// Additional Tank heroes
	"ramattra": func() any { return RamattraMetrics{} },

	// Additional Support heroes
	"moira":  func() any { return MoiraMetrics{} },
	"juno":   func() any { return JunoMetrics{} },
	"wuyang": func() any { return WuyangMetrics{} },
	"freja":  func() any { return FrejaMetrics{} },
}

// GetHeroMetrics returns MetricDef map for a specific hero (defaults to PC QuickPlay).
func GetHeroMetrics(heroID string) map[string]MetricDef {
	return GetHeroMetricsForPlatform(heroID, PlatformPC, GameModeQuickPlay)
}

// GetHeroMetricsForPlatform returns MetricDef map for a specific hero with platform and gamemode context.
func GetHeroMetricsForPlatform(heroID string, platform Platform, gameMode GameMode) map[string]MetricDef {
	factory, exists := HeroMetricsRegistry[heroID]
	if !exists {
		// Fallback to common metrics only with platform context
		return generateCommonMetricsWithContext(platform, gameMode)
	}

	heroStruct := factory()
	heroSpecific := GenerateMetricDefsWithContext(heroStruct, platform, gameMode)

	// Merge with common metrics from embedded struct
	commonStruct := CommonMetrics{}
	commonMetrics := GenerateMetricDefsWithContext(commonStruct, platform, gameMode)

	// Combine both maps
	result := make(map[string]MetricDef)
	for k, v := range commonMetrics {
		result[k] = v
	}

	for k, v := range heroSpecific {
		result[k] = v
	}

	return result
}

// generateCommonMetricsWithContext generates common metrics with platform context.
func generateCommonMetricsWithContext(platform Platform, gameMode GameMode) map[string]MetricDef {
	commonStruct := CommonMetrics{}

	return GenerateMetricDefsWithContext(commonStruct, platform, gameMode)
}

// generatePlatformSelector creates platform/gamemode aware CSS selectors.
func generatePlatformSelector(baseSelector string, platform Platform, gameMode GameMode) string {
	if baseSelector == "" {
		return ""
	}

	// Platform selector wrapper
	var platformWrapper string

	switch platform {
	case PlatformPC:
		platformWrapper = MouseKeyboardViewActiveSelector
	case PlatformConsole:
		platformWrapper = ".controller-view.is-active"
	default:
		platformWrapper = MouseKeyboardViewActiveSelector // Default to PC
	}

	// GameMode selector wrapper
	var gameModeWrapper string

	switch gameMode {
	case GameModeQuickPlay:
		gameModeWrapper = QuickPlayViewActiveSelector
	case GameModeCompetitive:
		gameModeWrapper = ".competitive-view.is-active"
	default:
		gameModeWrapper = QuickPlayViewActiveSelector // Default to QuickPlay
	}

	// Combine platform + gamemode + base selector
	return platformWrapper + " " + gameModeWrapper + " " + baseSelector
}

// IllariMetrics defines Illari specific metrics.
type IllariMetrics struct {
	CommonMetrics // Embedded common metrics

	SolarRifleHealing   int64 `ow:"solar_rifle_healing" prometheus:"ow_hero_solar_rifle_healing_total" help:"Total healing done with solar rifle" path:"[data-stat='solar_rifle_healing']" type:"number"`
	HealingPylonHealing int64 `ow:"healing_pylon_healing" prometheus:"ow_hero_healing_pylon_healing_total" help:"Total healing provided by healing pylon" path:"[data-stat='healing_pylon_healing']" type:"number"`
	CaptiveSunKills     int   `ow:"captive_sun_kills" prometheus:"ow_hero_captive_sun_kills_total" help:"Eliminations with captive sun ultimate" path:"[data-stat='captive_sun_kills']" type:"number"`
	OutburstKills       int   `ow:"outburst_kills" prometheus:"ow_hero_outburst_kills_total" help:"Eliminations with outburst ability" path:"[data-stat='outburst_kills']" type:"number"`
}

// CassidyMetrics defines Cassidy specific metrics.
type CassidyMetrics struct {
	CommonMetrics // Embedded common metrics

	PeacekeeperAccuracy float64 `ow:"peacekeeper_accuracy" prometheus:"ow_hero_peacekeeper_accuracy_percent" help:"Peacekeeper weapon accuracy percentage" path:"[data-stat='peacekeeper_accuracy']" type:"percentage"`
	FlashbangEnemies    int     `ow:"flashbang_enemies" prometheus:"ow_hero_flashbang_enemies_total" help:"Total enemies stunned with flashbang" path:"[data-stat='flashbang_enemies']" type:"number"`
	CombatRollKills     int     `ow:"combat_roll_kills" prometheus:"ow_hero_combat_roll_kills_total" help:"Eliminations after using combat roll" path:"[data-stat='combat_roll_kills']" type:"number"`
	DeadeyeKills        int     `ow:"deadeye_kills" prometheus:"ow_hero_deadeye_kills_total" help:"Eliminations with deadeye ultimate" path:"[data-stat='deadeye_kills']" type:"number"`
}

// LifeweaverMetrics defines Lifeweaver specific metrics.
type LifeweaverMetrics struct {
	CommonMetrics // Embedded common metrics

	HealingBlossomHealing int64         `ow:"healing_blossom_healing" prometheus:"ow_hero_healing_blossom_healing_total" help:"Total healing done with healing blossom" path:"[data-stat='healing_blossom_healing']" type:"number"`
	LifeGripSaves         int           `ow:"life_grip_saves" prometheus:"ow_hero_life_grip_saves_total" help:"Teammates saved with life grip" path:"[data-stat='life_grip_saves']" type:"number"`
	PetalPlatformUptime   time.Duration `ow:"petal_platform_uptime" prometheus:"ow_hero_petal_platform_uptime_seconds" help:"Total uptime of petal platforms" path:"[data-stat='petal_platform_uptime']" type:"duration"`
	TreeOfLifeHealing     int64         `ow:"tree_of_life_healing" prometheus:"ow_hero_tree_of_life_healing_total" help:"Healing provided by tree of life ultimate" path:"[data-stat='tree_of_life_healing']" type:"number"`
}

// TracerMetrics defines Tracer specific metrics.
type TracerMetrics struct {
	CommonMetrics // Embedded common metrics

	PulseGunsAccuracy float64 `ow:"pulse_guns_accuracy" prometheus:"ow_hero_pulse_guns_accuracy_percent" help:"Pulse guns weapon accuracy percentage" path:"[data-stat='pulse_guns_accuracy']" type:"percentage"`
	BlinkDistance     float64 `ow:"blink_distance" prometheus:"ow_hero_blink_distance_meters" help:"Total distance traveled with blink" path:"[data-stat='blink_distance']" type:"number"`
	RecallHealing     int64   `ow:"recall_healing" prometheus:"ow_hero_recall_healing_total" help:"Health recovered using recall" path:"[data-stat='recall_healing']" type:"number"`
	PulseBombKills    int     `ow:"pulse_bomb_kills" prometheus:"ow_hero_pulse_bomb_kills_total" help:"Eliminations with pulse bomb ultimate" path:"[data-stat='pulse_bomb_kills']" type:"number"`
}

// KirikoMetrics defines Kiriko specific metrics.
type KirikoMetrics struct {
	CommonMetrics // Embedded common metrics

	HealingOfuudaHealing int64 `ow:"healing_ofuuda_healing" prometheus:"ow_hero_healing_ofuuda_healing_total" help:"Total healing done with healing ofuuda" path:"[data-stat='healing_ofuuda_healing']" type:"number"`
	KunaiCriticalHits    int   `ow:"kunai_critical_hits" prometheus:"ow_hero_kunai_critical_hits_total" help:"Critical hits with kunai" path:"[data-stat='kunai_critical_hits']" type:"number"`
	SwiftStepTeleports   int   `ow:"swift_step_teleports" prometheus:"ow_hero_swift_step_teleports_total" help:"Number of swift step teleports used" path:"[data-stat='swift_step_teleports']" type:"number"`
	KitsuneFinalBlows    int   `ow:"kitsune_final_blows" prometheus:"ow_hero_kitsune_final_blows_total" help:"Final blows during kitsune rush ultimate" path:"[data-stat='kitsune_final_blows']" type:"number"`
}

// PharahMetrics defines Pharah specific metrics.
type PharahMetrics struct {
	CommonMetrics // Embedded common metrics

	RocketLauncherAccuracy float64       `ow:"rocket_launcher_accuracy" prometheus:"ow_hero_rocket_launcher_accuracy_percent" help:"Rocket launcher weapon accuracy percentage" path:"[data-stat='rocket_launcher_accuracy']" type:"percentage"`
	ConcussiveBlastKills   int           `ow:"concussive_blast_kills" prometheus:"ow_hero_concussive_blast_kills_total" help:"Environmental kills with concussive blast" path:"[data-stat='concussive_blast_kills']" type:"number"`
	BarrageKills           int           `ow:"barrage_kills" prometheus:"ow_hero_barrage_kills_total" help:"Eliminations with barrage ultimate" path:"[data-stat='barrage_kills']" type:"number"`
	AirborneTime           time.Duration `ow:"airborne_time" prometheus:"ow_hero_airborne_time_seconds" help:"Total time spent airborne" path:"[data-stat='airborne_time']" type:"duration"`
}

// WinstonMetrics defines Winston specific metrics.
type WinstonMetrics struct {
	CommonMetrics // Embedded common metrics

	TeslaCcnnonKills       int           `ow:"tesla_cannon_kills" prometheus:"ow_hero_tesla_cannon_kills_total" help:"Eliminations with tesla cannon" path:"[data-stat='tesla_cannon_kills']" type:"number"`
	JumpPackKills          int           `ow:"jump_pack_kills" prometheus:"ow_hero_jump_pack_kills_total" help:"Eliminations with jump pack" path:"[data-stat='jump_pack_kills']" type:"number"`
	BarrierProjectorUptime time.Duration `ow:"barrier_projector_uptime" prometheus:"ow_hero_barrier_projector_uptime_seconds" help:"Total uptime of barrier projector" path:"[data-stat='barrier_projector_uptime']" type:"duration"`
	PrimalRageKills        int           `ow:"primal_rage_kills" prometheus:"ow_hero_primal_rage_kills_total" help:"Eliminations during primal rage ultimate" path:"[data-stat='primal_rage_kills']" type:"number"`
}

// RoadhogMetrics defines Roadhog specific metrics.
type RoadhogMetrics struct {
	CommonMetrics // Embedded common metrics

	ScrapGunAccuracy    float64 `ow:"scrap_gun_accuracy" prometheus:"ow_hero_scrap_gun_accuracy_percent" help:"Scrap gun weapon accuracy percentage" path:"[data-stat='scrap_gun_accuracy']" type:"percentage"`
	ChainHookAccuracy   float64 `ow:"chain_hook_accuracy" prometheus:"ow_hero_chain_hook_accuracy_percent" help:"Chain hook accuracy percentage" path:"[data-stat='chain_hook_accuracy']" type:"percentage"`
	ChainHookKills      int     `ow:"chain_hook_kills" prometheus:"ow_hero_chain_hook_kills_total" help:"Eliminations after chain hook" path:"[data-stat='chain_hook_kills']" type:"number"`
	TakeABreatheHealing int64   `ow:"take_a_breathe_healing" prometheus:"ow_hero_take_a_breathe_healing_total" help:"Health recovered with take a breathe" path:"[data-stat='take_a_breathe_healing']" type:"number"`
	WholePigKills       int     `ow:"whole_pig_kills" prometheus:"ow_hero_whole_pig_kills_total" help:"Eliminations with whole hog ultimate" path:"[data-stat='whole_pig_kills']" type:"number"`
}

// SojournMetrics defines Sojourn specific metrics.
type SojournMetrics struct {
	CommonMetrics // Embedded common metrics

	RailgunAccuracy     float64 `ow:"railgun_accuracy" prometheus:"ow_hero_railgun_accuracy_percent" help:"Railgun weapon accuracy percentage" path:"[data-stat='railgun_accuracy']" type:"percentage"`
	RailgunCriticalHits int     `ow:"railgun_critical_hits" prometheus:"ow_hero_railgun_critical_hits_total" help:"Critical hits with railgun" path:"[data-stat='railgun_critical_hits']" type:"number"`
	PowerSlideKills     int     `ow:"power_slide_kills" prometheus:"ow_hero_power_slide_kills_total" help:"Eliminations using power slide" path:"[data-stat='power_slide_kills']" type:"number"`
	OverclockKills      int     `ow:"overclock_kills" prometheus:"ow_hero_overclock_kills_total" help:"Eliminations during overclock ultimate" path:"[data-stat='overclock_kills']" type:"number"`
}

// BaptisteMetrics defines Baptiste specific metrics.
type BaptisteMetrics struct {
	CommonMetrics // Embedded common metrics

	BioticLauncherHealing     int64 `ow:"biotic_launcher_healing" prometheus:"ow_hero_biotic_launcher_healing_total" help:"Total healing done with biotic launcher" path:"[data-stat='biotic_launcher_healing']" type:"number"`
	RegenerativeBurstHealing  int64 `ow:"regenerative_burst_healing" prometheus:"ow_hero_regenerative_burst_healing_total" help:"Healing provided by regenerative burst" path:"[data-stat='regenerative_burst_healing']" type:"number"`
	ImmortalityFieldSaves     int   `ow:"immortality_field_saves" prometheus:"ow_hero_immortality_field_saves_total" help:"Teammates saved with immortality field" path:"[data-stat='immortality_field_saves']" type:"number"`
	AmplificationMatrixDamage int64 `ow:"amplification_matrix_damage" prometheus:"ow_hero_amplification_matrix_damage_total" help:"Damage amplified by amplification matrix ultimate" path:"[data-stat='amplification_matrix_damage']" type:"number"`
}

// MeiMetrics defines Mei specific metrics.
type MeiMetrics struct {
	CommonMetrics // Embedded common metrics

	EndothermlcBlasterAccuracy float64       `ow:"endothermic_blaster_accuracy" prometheus:"ow_hero_endothermic_blaster_accuracy_percent" help:"Endothermic blaster weapon accuracy percentage" path:"[data-stat='endothermic_blaster_accuracy']" type:"percentage"`
	EnemiesFrozen              int           `ow:"enemies_frozen" prometheus:"ow_hero_enemies_frozen_total" help:"Total enemies frozen" path:"[data-stat='enemies_frozen']" type:"number"`
	IceWallUptime              time.Duration `ow:"ice_wall_uptime" prometheus:"ow_hero_ice_wall_uptime_seconds" help:"Total uptime of ice walls" path:"[data-stat='ice_wall_uptime']" type:"duration"`
	BlizzardKills              int           `ow:"blizzard_kills" prometheus:"ow_hero_blizzard_kills_total" help:"Eliminations with blizzard ultimate" path:"[data-stat='blizzard_kills']" type:"number"`
}

// ZaryaMetrics defines Zarya specific metrics.
type ZaryaMetrics struct {
	CommonMetrics // Embedded common metrics

	ParticannonKills              int   `ow:"particle_cannon_kills" prometheus:"ow_hero_particle_cannon_kills_total" help:"Eliminations with particle cannon" path:"[data-stat='particle_cannon_kills']" type:"number"`
	ParticleBarrierDamageAbsorbed int64 `ow:"particle_barrier_absorbed" prometheus:"ow_hero_particle_barrier_absorbed_total" help:"Damage absorbed by particle barriers" path:"[data-stat='particle_barrier_absorbed']" type:"number"`
	ProjectedBarrierSaves         int   `ow:"projected_barrier_saves" prometheus:"ow_hero_projected_barrier_saves_total" help:"Teammates saved with projected barrier" path:"[data-stat='projected_barrier_saves']" type:"number"`
	GravitonSurgeKills            int   `ow:"graviton_surge_kills" prometheus:"ow_hero_graviton_surge_kills_total" help:"Eliminations with graviton surge ultimate" path:"[data-stat='graviton_surge_kills']" type:"number"`
	HighEnergyKills               int   `ow:"high_energy_kills" prometheus:"ow_hero_high_energy_kills_total" help:"Eliminations while at high energy" path:"[data-stat='high_energy_kills']" type:"number"`
}

// JunkratMetrics defines Junkrat specific metrics.
type JunkratMetrics struct {
	CommonMetrics // Embedded common metrics

	FragLauncherAccuracy float64 `ow:"frag_launcher_accuracy" prometheus:"ow_hero_frag_launcher_accuracy_percent" help:"Frag launcher weapon accuracy percentage" path:"[data-stat='frag_launcher_accuracy']" type:"percentage"`
	ConcussionMineKills  int     `ow:"concussion_mine_kills" prometheus:"ow_hero_concussion_mine_kills_total" help:"Eliminations with concussion mine" path:"[data-stat='concussion_mine_kills']" type:"number"`
	SteelTrapKills       int     `ow:"steel_trap_kills" prometheus:"ow_hero_steel_trap_kills_total" help:"Eliminations with steel trap" path:"[data-stat='steel_trap_kills']" type:"number"`
	RipTireKills         int     `ow:"rip_tire_kills" prometheus:"ow_hero_rip_tire_kills_total" help:"Eliminations with rip-tire ultimate" path:"[data-stat='rip_tire_kills']" type:"number"`
	EnemiesTrapped       int     `ow:"enemies_trapped" prometheus:"ow_hero_enemies_trapped_total" help:"Enemies caught in steel trap" path:"[data-stat='enemies_trapped']" type:"number"`
}

// LucioMetrics defines Lúcio specific metrics.
type LucioMetrics struct {
	CommonMetrics // Embedded common metrics

	SonicAmplifierAccuracy float64       `ow:"sonic_amplifier_accuracy" prometheus:"ow_hero_sonic_amplifier_accuracy_percent" help:"Sonic amplifier weapon accuracy percentage" path:"[data-stat='sonic_amplifier_accuracy']" type:"percentage"`
	SoundBoopKills         int           `ow:"sound_boop_kills" prometheus:"ow_hero_sound_boop_kills_total" help:"Environmental kills with sound wave" path:"[data-stat='sound_boop_kills']" type:"number"`
	WallRideTime           time.Duration `ow:"wall_ride_time" prometheus:"ow_hero_wall_ride_time_seconds" help:"Total time spent wall riding" path:"[data-stat='wall_ride_time']" type:"duration"`
	SoundBarrierProvided   int64         `ow:"sound_barrier_provided" prometheus:"ow_hero_sound_barrier_provided_total" help:"Shield health provided by sound barrier ultimate" path:"[data-stat='sound_barrier_provided']" type:"number"`
}

// ReaperMetrics defines Reaper specific metrics.
type ReaperMetrics struct {
	CommonMetrics // Embedded common metrics

	HellfireshotgunsAccuracy float64 `ow:"hellfire_shotguns_accuracy" prometheus:"ow_hero_hellfire_shotguns_accuracy_percent" help:"Hellfire shotguns weapon accuracy percentage" path:"[data-stat='hellfire_shotguns_accuracy']" type:"percentage"`
	WraithFormDistance       float64 `ow:"wraith_form_distance" prometheus:"ow_hero_wraith_form_distance_meters" help:"Distance traveled in wraith form" path:"[data-stat='wraith_form_distance']" type:"number"`
	ShadowStepTeleports      int     `ow:"shadow_step_teleports" prometheus:"ow_hero_shadow_step_teleports_total" help:"Number of shadow step teleports" path:"[data-stat='shadow_step_teleports']" type:"number"`
	DeathBlossomKills        int     `ow:"death_blossom_kills" prometheus:"ow_hero_death_blossom_kills_total" help:"Eliminations with death blossom ultimate" path:"[data-stat='death_blossom_kills']" type:"number"`
}

// ZenyattaMetrics defines Zenyatta specific metrics.
type ZenyattaMetrics struct {
	CommonMetrics // Embedded common metrics

	OrbOfDestructionAccuracy float64 `ow:"orb_destruction_accuracy" prometheus:"ow_hero_orb_destruction_accuracy_percent" help:"Orb of destruction weapon accuracy percentage" path:"[data-stat='orb_destruction_accuracy']" type:"percentage"`
	OrbOfDiscordAssists      int     `ow:"orb_discord_assists" prometheus:"ow_hero_orb_discord_assists_total" help:"Eliminations assisted by orb of discord" path:"[data-stat='orb_discord_assists']" type:"number"`
	OrbOfHarmonyHealing      int64   `ow:"orb_harmony_healing" prometheus:"ow_hero_orb_harmony_healing_total" help:"Healing provided by orb of harmony" path:"[data-stat='orb_harmony_healing']" type:"number"`
	TranscendenceHealing     int64   `ow:"transcendence_healing" prometheus:"ow_hero_transcendence_healing_total" help:"Healing provided during transcendence ultimate" path:"[data-stat='transcendence_healing']" type:"number"`
}

// MaugaMetrics defines Mauga specific metrics.
type MaugaMetrics struct {
	CommonMetrics // Embedded common metrics

	IncendiaryChaingunDamage int64 `ow:"incendiary_chaingun_damage" prometheus:"ow_hero_incendiary_chaingun_damage_total" help:"Damage dealt with incendiary chaingun" path:"[data-stat='incendiary_chaingun_damage']" type:"number"`
	VolatileChaingunDamage   int64 `ow:"volatile_chaingun_damage" prometheus:"ow_hero_volatile_chaingun_damage_total" help:"Damage dealt with volatile chaingun" path:"[data-stat='volatile_chaingun_damage']" type:"number"`
	OverrunKills             int   `ow:"overrun_kills" prometheus:"ow_hero_overrun_kills_total" help:"Eliminations with overrun charge" path:"[data-stat='overrun_kills']" type:"number"`
	CageFightKills           int   `ow:"cage_fight_kills" prometheus:"ow_hero_cage_fight_kills_total" help:"Eliminations during cage fight ultimate" path:"[data-stat='cage_fight_kills']" type:"number"`
}

// BrigitteMetrics defines Brigitte specific metrics.
type BrigitteMetrics struct {
	CommonMetrics // Embedded common metrics

	RocketFlailAccuracy        float64 `ow:"rocket_flail_accuracy" prometheus:"ow_hero_rocket_flail_accuracy_percent" help:"Rocket flail weapon accuracy percentage" path:"[data-stat='rocket_flail_accuracy']" type:"percentage"`
	WhipShotKills              int     `ow:"whip_shot_kills" prometheus:"ow_hero_whip_shot_kills_total" help:"Environmental kills with whip shot" path:"[data-stat='whip_shot_kills']" type:"number"`
	RepairPackHealing          int64   `ow:"repair_pack_healing" prometheus:"ow_hero_repair_pack_healing_total" help:"Healing provided by repair pack" path:"[data-stat='repair_pack_healing']" type:"number"`
	BarrierShieldDamageBlocked int64   `ow:"barrier_shield_blocked" prometheus:"ow_hero_barrier_shield_blocked_total" help:"Damage blocked by barrier shield" path:"[data-stat='barrier_shield_blocked']" type:"number"`
	RallyShieldProvided        int64   `ow:"rally_shield_provided" prometheus:"ow_hero_rally_shield_provided_total" help:"Shield health provided by rally ultimate" path:"[data-stat='rally_shield_provided']" type:"number"`
}

// HazardMetrics defines Hazard specific metrics.
type HazardMetrics struct {
	CommonMetrics // Embedded common metrics

	SpikeTrapKills    int           `ow:"spike_trap_kills" prometheus:"ow_hero_spike_trap_kills_total" help:"Eliminations with spike trap" path:"[data-stat='spike_trap_kills']" type:"number"`
	ViolentLeapKills  int           `ow:"violent_leap_kills" prometheus:"ow_hero_violent_leap_kills_total" help:"Eliminations with violent leap" path:"[data-stat='violent_leap_kills']" type:"number"`
	DowntimeUptime    time.Duration `ow:"downtime_uptime" prometheus:"ow_hero_downtime_uptime_seconds" help:"Total uptime of downtime walls" path:"[data-stat='downtime_uptime']" type:"duration"`
	VanadiumRageKills int           `ow:"vanadium_rage_kills" prometheus:"ow_hero_vanadium_rage_kills_total" help:"Eliminations during vanadium rage ultimate" path:"[data-stat='vanadium_rage_kills']" type:"number"`
}

// JunkerQueenMetrics defines Junker Queen specific metrics.
type JunkerQueenMetrics struct {
	CommonMetrics // Embedded common metrics

	ScattergunAccuracy     float64 `ow:"scattergun_accuracy" prometheus:"ow_hero_scattergun_accuracy_percent" help:"Scattergun weapon accuracy percentage" path:"[data-stat='scattergun_accuracy']" type:"percentage"`
	JaggedBladeKills       int     `ow:"jagged_blade_kills" prometheus:"ow_hero_jagged_blade_kills_total" help:"Eliminations with jagged blade" path:"[data-stat='jagged_blade_kills']" type:"number"`
	CommandingShoutHealing int64   `ow:"commanding_shout_healing" prometheus:"ow_hero_commanding_shout_healing_total" help:"Healing provided by commanding shout" path:"[data-stat='commanding_shout_healing']" type:"number"`
	RampageKills           int     `ow:"rampage_kills" prometheus:"ow_hero_rampage_kills_total" help:"Eliminations with rampage ultimate" path:"[data-stat='rampage_kills']" type:"number"`
}

// HanzoMetrics defines Hanzo specific metrics.
type HanzoMetrics struct {
	CommonMetrics // Embedded common metrics

	StormBowAccuracy  float64 `ow:"storm_bow_accuracy" prometheus:"ow_hero_storm_bow_accuracy_percent" help:"Storm bow weapon accuracy percentage" path:"[data-stat='storm_bow_accuracy']" type:"percentage"`
	SonicArrowAssists int     `ow:"sonic_arrow_assists" prometheus:"ow_hero_sonic_arrow_assists_total" help:"Eliminations assisted by sonic arrow" path:"[data-stat='sonic_arrow_assists']" type:"number"`
	ScatterArrowKills int     `ow:"scatter_arrow_kills" prometheus:"ow_hero_scatter_arrow_kills_total" help:"Eliminations with scatter arrow" path:"[data-stat='scatter_arrow_kills']" type:"number"`
	DragonstrikeKills int     `ow:"dragonstrike_kills" prometheus:"ow_hero_dragonstrike_kills_total" help:"Eliminations with dragonstrike ultimate" path:"[data-stat='dragonstrike_kills']" type:"number"`
}

// DoomfistMetrics defines Doomfist specific metrics.
type DoomfistMetrics struct {
	CommonMetrics // Embedded common metrics

	HandCannonAccuracy float64 `ow:"hand_cannon_accuracy" prometheus:"ow_hero_hand_cannon_accuracy_percent" help:"Hand cannon weapon accuracy percentage" path:"[data-stat='hand_cannon_accuracy']" type:"percentage"`
	SeismicSlamKills   int     `ow:"seismic_slam_kills" prometheus:"ow_hero_seismic_slam_kills_total" help:"Eliminations with seismic slam" path:"[data-stat='seismic_slam_kills']" type:"number"`
	RocketPunchKills   int     `ow:"rocket_punch_kills" prometheus:"ow_hero_rocket_punch_kills_total" help:"Eliminations with rocket punch" path:"[data-stat='rocket_punch_kills']" type:"number"`
	MeteorStrikeKills  int     `ow:"meteor_strike_kills" prometheus:"ow_hero_meteor_strike_kills_total" help:"Eliminations with meteor strike ultimate" path:"[data-stat='meteor_strike_kills']" type:"number"`
}

// MoiraMetrics defines Moira specific metrics.
type MoiraMetrics struct {
	CommonMetrics // Embedded common metrics

	BiotiGraspAccuracy float64 `ow:"biotic_grasp_accuracy" prometheus:"ow_hero_biotic_grasp_accuracy_percent" help:"Biotic grasp weapon accuracy percentage" path:"[data-stat='biotic_grasp_accuracy']" type:"percentage"`
	CoalescenceKills   int     `ow:"coalescence_kills" prometheus:"ow_hero_coalescence_kills_total" help:"Eliminations with coalescence ultimate" path:"[data-stat='coalescence_kills']" type:"number"`
	BiotiOrbKills      int     `ow:"biotic_orb_kills" prometheus:"ow_hero_biotic_orb_kills_total" help:"Eliminations with biotic orb" path:"[data-stat='biotic_orb_kills']" type:"number"`
	SelfHealing        int64   `ow:"self_healing" prometheus:"ow_hero_self_healing_total" help:"Self healing done" path:"[data-stat='self_healing']" type:"number"`
}

// OrisaMetrics defines Orisa specific metrics.
type OrisaMetrics struct {
	CommonMetrics // Embedded common metrics

	FusionDriverAccuracy float64 `ow:"fusion_driver_accuracy" prometheus:"ow_hero_fusion_driver_accuracy_percent" help:"Fusion driver weapon accuracy percentage" path:"[data-stat='fusion_driver_accuracy']" type:"percentage"`
	EnergyJavelinKills   int     `ow:"energy_javelin_kills" prometheus:"ow_hero_energy_javelin_kills_total" help:"Eliminations with energy javelin" path:"[data-stat='energy_javelin_kills']" type:"number"`
	TerraForceKills      int     `ow:"terra_force_kills" prometheus:"ow_hero_terra_force_kills_total" help:"Eliminations with terra force ultimate" path:"[data-stat='terra_force_kills']" type:"number"`
	DamageAmplified      int64   `ow:"damage_amplified" prometheus:"ow_hero_damage_amplified_total" help:"Damage amplified for teammates" path:"[data-stat='damage_amplified']" type:"number"`
}

// SigmaMetrics defines Sigma specific metrics.
type SigmaMetrics struct {
	CommonMetrics // Embedded common metrics

	HyperSpheresAccuracy float64 `ow:"hyper_spheres_accuracy" prometheus:"ow_hero_hyper_spheres_accuracy_percent" help:"Hyper spheres weapon accuracy percentage" path:"[data-stat='hyper_spheres_accuracy']" type:"percentage"`
	AccretionKills       int     `ow:"accretion_kills" prometheus:"ow_hero_accretion_kills_total" help:"Eliminations with accretion" path:"[data-stat='accretion_kills']" type:"number"`
	GraviticFluxKills    int     `ow:"gravitic_flux_kills" prometheus:"ow_hero_gravitic_flux_kills_total" help:"Eliminations with gravitic flux ultimate" path:"[data-stat='gravitic_flux_kills']" type:"number"`
	DamageAbsorbed       int64   `ow:"damage_absorbed" prometheus:"ow_hero_damage_absorbed_total" help:"Damage absorbed by experimental barrier" path:"[data-stat='damage_absorbed']" type:"number"`
}

// SombraMetrics defines Sombra specific metrics.
type SombraMetrics struct {
	CommonMetrics // Embedded common metrics

	MachineGunAccuracy float64 `ow:"machine_gun_accuracy" prometheus:"ow_hero_machine_gun_accuracy_percent" help:"Machine gun weapon accuracy percentage" path:"[data-stat='machine_gun_accuracy']" type:"percentage"`
	EnemiesHacked      int     `ow:"enemies_hacked" prometheus:"ow_hero_enemies_hacked_total" help:"Total enemies hacked" path:"[data-stat='enemies_hacked']" type:"number"`
	EnemiesEMPd        int     `ow:"enemies_empd" prometheus:"ow_hero_enemies_empd_total" help:"Total enemies hit with EMP ultimate" path:"[data-stat='enemies_empd']" type:"number"`
	HealthPacksHacked  int     `ow:"health_packs_hacked" prometheus:"ow_hero_health_packs_hacked_total" help:"Total health packs hacked" path:"[data-stat='health_packs_hacked']" type:"number"`
}

// SymmetraMetrics defines Symmetra specific metrics.
type SymmetraMetrics struct {
	CommonMetrics // Embedded common metrics

	PhotonProjectorAccuracy float64 `ow:"photon_projector_accuracy" prometheus:"ow_hero_photon_projector_accuracy_percent" help:"Photon projector weapon accuracy percentage" path:"[data-stat='photon_projector_accuracy']" type:"percentage"`
	SentryTurretKills       int     `ow:"sentry_turret_kills" prometheus:"ow_hero_sentry_turret_kills_total" help:"Eliminations with sentry turrets" path:"[data-stat='sentry_turret_kills']" type:"number"`
	TeleporterPadsSummoned  int     `ow:"teleporter_pads_summoned" prometheus:"ow_hero_teleporter_pads_summoned_total" help:"Total teleporter pads summoned" path:"[data-stat='teleporter_pads_summoned']" type:"number"`
	PlayersTeleported       int     `ow:"players_teleported" prometheus:"ow_hero_players_teleported_total" help:"Total players teleported" path:"[data-stat='players_teleported']" type:"number"`
}

// WreckingBallMetrics defines Wrecking Ball specific metrics.
type WreckingBallMetrics struct {
	CommonMetrics // Embedded common metrics

	QuadCannonAccuracy float64 `ow:"quad_cannon_accuracy" prometheus:"ow_hero_quad_cannon_accuracy_percent" help:"Quad cannon weapon accuracy percentage" path:"[data-stat='quad_cannon_accuracy']" type:"percentage"`
	PiledriveKills     int     `ow:"piledrive_kills" prometheus:"ow_hero_piledrive_kills_total" help:"Eliminations with piledrive" path:"[data-stat='piledrive_kills']" type:"number"`
	MinefieldKills     int     `ow:"minefield_kills" prometheus:"ow_hero_minefield_kills_total" help:"Eliminations with minefield ultimate" path:"[data-stat='minefield_kills']" type:"number"`
	PlayersKnockedBack int     `ow:"players_knocked_back" prometheus:"ow_hero_players_knocked_back_total" help:"Total players knocked back" path:"[data-stat='players_knocked_back']" type:"number"`
}

// BastionMetrics defines Bastion specific metrics.
type BastionMetrics struct {
	CommonMetrics // Embedded common metrics

	ConfigurationAssaultAccuracy float64 `ow:"configuration_assault_accuracy" prometheus:"ow_hero_configuration_assault_accuracy_percent" help:"Configuration assault weapon accuracy percentage" path:"[data-stat='configuration_assault_accuracy']" type:"percentage"`
	ConfigurationReconAccuracy   float64 `ow:"configuration_recon_accuracy" prometheus:"ow_hero_configuration_recon_accuracy_percent" help:"Configuration recon weapon accuracy percentage" path:"[data-stat='configuration_recon_accuracy']" type:"percentage"`
	SelfRepairUsed               int     `ow:"self_repair_used" prometheus:"ow_hero_self_repair_used_total" help:"Total times self repair was used" path:"[data-stat='self_repair_used']" type:"number"`
	ConfigurationArtilleryKills  int     `ow:"configuration_artillery_kills" prometheus:"ow_hero_configuration_artillery_kills_total" help:"Eliminations with configuration artillery ultimate" path:"[data-stat='configuration_artillery_kills']" type:"number"`
}

// AsheMetrics defines Ashe specific metrics.
type AsheMetrics struct {
	CommonMetrics // Embedded common metrics

	ViperAccuracy float64 `ow:"viper_accuracy" prometheus:"ow_hero_viper_accuracy_percent" help:"Viper rifle weapon accuracy percentage" path:"[data-stat='viper_accuracy']" type:"percentage"`
	DynamiteKills int     `ow:"dynamite_kills" prometheus:"ow_hero_dynamite_kills_total" help:"Eliminations with dynamite" path:"[data-stat='dynamite_kills']" type:"number"`
	CoachGunKills int     `ow:"coach_gun_kills" prometheus:"ow_hero_coach_gun_kills_total" help:"Eliminations with coach gun" path:"[data-stat='coach_gun_kills']" type:"number"`
	BOBKills      int     `ow:"bob_kills" prometheus:"ow_hero_bob_kills_total" help:"Eliminations with BOB ultimate" path:"[data-stat='bob_kills']" type:"number"`
}

// EchoMetrics defines Echo specific metrics.
type EchoMetrics struct {
	CommonMetrics // Embedded common metrics

	TriShotAccuracy        float64 `ow:"tri_shot_accuracy" prometheus:"ow_hero_tri_shot_accuracy_percent" help:"Tri-shot weapon accuracy percentage" path:"[data-stat='tri_shot_accuracy']" type:"percentage"`
	StickyBombKills        int     `ow:"sticky_bomb_kills" prometheus:"ow_hero_sticky_bomb_kills_total" help:"Eliminations with sticky bombs" path:"[data-stat='sticky_bomb_kills']" type:"number"`
	FlightTimeUsed         int     `ow:"flight_time_used" prometheus:"ow_hero_flight_time_used_total" help:"Total flight time used" path:"[data-stat='flight_time_used']" type:"number"`
	DuplicateUltimateKills int     `ow:"duplicate_ultimate_kills" prometheus:"ow_hero_duplicate_ultimate_kills_total" help:"Eliminations with duplicated ultimates" path:"[data-stat='duplicate_ultimate_kills']" type:"number"`
}

// VentureMetrics defines Venture specific metrics.
type VentureMetrics struct {
	CommonMetrics // Embedded common metrics

	SmartExcavatorAccuracy float64 `ow:"smart_excavator_accuracy" prometheus:"ow_hero_smart_excavator_accuracy_percent" help:"Smart excavator weapon accuracy percentage" path:"[data-stat='smart_excavator_accuracy']" type:"percentage"`
	BurrowKills            int     `ow:"burrow_kills" prometheus:"ow_hero_burrow_kills_total" help:"Eliminations with burrow" path:"[data-stat='burrow_kills']" type:"number"`
	DrillDashKills         int     `ow:"drill_dash_kills" prometheus:"ow_hero_drill_dash_kills_total" help:"Eliminations with drill dash" path:"[data-stat='drill_dash_kills']" type:"number"`
	TectonicShockKills     int     `ow:"tectonic_shock_kills" prometheus:"ow_hero_tectonic_shock_kills_total" help:"Eliminations with tectonic shock ultimate" path:"[data-stat='tectonic_shock_kills']" type:"number"`
}

// RamattraMetrics defines Ramattra specific metrics.
type RamattraMetrics struct {
	CommonMetrics // Embedded common metrics

	VoidAcceleratorAccuracy  float64 `ow:"void_accelerator_accuracy" prometheus:"ow_hero_void_accelerator_accuracy_percent" help:"Void accelerator weapon accuracy percentage" path:"[data-stat='void_accelerator_accuracy']" type:"percentage"`
	VoidBarrierDamageBlocked int64   `ow:"void_barrier_damage_blocked" prometheus:"ow_hero_void_barrier_damage_blocked_total" help:"Damage blocked by void barrier" path:"[data-stat='void_barrier_damage_blocked']" type:"number"`
	RavenousVortexKills      int     `ow:"ravenous_vortex_kills" prometheus:"ow_hero_ravenous_vortex_kills_total" help:"Eliminations with ravenous vortex" path:"[data-stat='ravenous_vortex_kills']" type:"number"`
	AnnihilationKills        int     `ow:"annihilation_kills" prometheus:"ow_hero_annihilation_kills_total" help:"Eliminations with annihilation ultimate" path:"[data-stat='annihilation_kills']" type:"number"`
}

// JunoMetrics defines Juno specific metrics.
type JunoMetrics struct {
	CommonMetrics // Embedded common metrics

	MediBlasterAccuracy float64 `ow:"medi_blaster_accuracy" prometheus:"ow_hero_medi_blaster_accuracy_percent" help:"Medi-blaster weapon accuracy percentage" path:"[data-stat='medi_blaster_accuracy']" type:"percentage"`
	PulsarTorpedoKills  int     `ow:"pulsar_torpedo_kills" prometheus:"ow_hero_pulsar_torpedo_kills_total" help:"Eliminations with pulsar torpedo" path:"[data-stat='pulsar_torpedo_kills']" type:"number"`
	GlideBoostUsed      int     `ow:"glide_boost_used" prometheus:"ow_hero_glide_boost_used_total" help:"Times glide boost was used" path:"[data-stat='glide_boost_used']" type:"number"`
	OrbitalRayKills     int     `ow:"orbital_ray_kills" prometheus:"ow_hero_orbital_ray_kills_total" help:"Eliminations with orbital ray ultimate" path:"[data-stat='orbital_ray_kills']" type:"number"`
}

// WuyangMetrics defines Wuyang specific metrics.
type WuyangMetrics struct {
	CommonMetrics // Embedded common metrics

	UnloadAccuracy      float64 `ow:"unload_accuracy" prometheus:"ow_hero_unload_accuracy_percent" help:"Unload weapon accuracy percentage" path:"[data-stat='unload_accuracy']" type:"percentage"`
	SteadyingKills      int     `ow:"steadying_kills" prometheus:"ow_hero_steadying_kills_total" help:"Eliminations with steadying" path:"[data-stat='steadying_kills']" type:"number"`
	CriticalMomentKills int     `ow:"critical_moment_kills" prometheus:"ow_hero_critical_moment_kills_total" help:"Eliminations with critical moment" path:"[data-stat='critical_moment_kills']" type:"number"`
	CriticalMomentUses  int     `ow:"critical_moment_uses" prometheus:"ow_hero_critical_moment_uses_total" help:"Times critical moment ability was used" path:"[data-stat='critical_moment_uses']" type:"number"`
}

// FrejaMetrics defines Freja specific metrics.
type FrejaMetrics struct {
	CommonMetrics // Embedded common metrics

	FrostBiteAccuracy    float64 `ow:"frost_bite_accuracy" prometheus:"ow_hero_frost_bite_accuracy_percent" help:"Frost bite weapon accuracy percentage" path:"[data-stat='frost_bite_accuracy']" type:"percentage"`
	IceWallUsed          int     `ow:"ice_wall_used" prometheus:"ow_hero_ice_wall_used_total" help:"Times ice wall was used" path:"[data-stat='ice_wall_used']" type:"number"`
	SlowDebuffApplied    int     `ow:"slow_debuff_applied" prometheus:"ow_hero_slow_debuff_applied_total" help:"Times slow debuff was applied" path:"[data-stat='slow_debuff_applied']" type:"number"`
	ArcticInversionKills int     `ow:"arctic_inversion_kills" prometheus:"ow_hero_arctic_inversion_kills_total" help:"Eliminations with arctic inversion ultimate" path:"[data-stat='arctic_inversion_kills']" type:"number"`
}

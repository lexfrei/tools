# Platform Analysis for Overwatch 2 Profiles

## Discovered Platform Structure

### Platform Tabs
```html
<blz-tab-control class="Profile-player--filter is-active" id="mouseKeyboardFilter">
    <blz-icon slot="icon" icon="bn-desktop-filled"></blz-icon>PC
</blz-tab-control>
<blz-tab-control class="Profile-player--filter" id="controllerFilter">
    <blz-icon slot="icon" icon="bn-console-filled"></blz-icon>Console
</blz-tab-control>
```

### Platform Views
```html
<!-- PC View (Active by default) -->
<div class="mouseKeyboard-view Profile-view is-active">
    <!-- Hero statistics for PC/Keyboard+Mouse -->
</div>

<!-- Console View (Hidden) -->
<div class="controller-view Profile-view">
    <!-- Hero statistics for Console/Controller -->
</div>
```

### Game Mode Tabs (within each platform)
```html
<blz-tab-control class="Profile-player--filter quickPlay-filter is-active">Quick Play</blz-tab-control>
<blz-tab-control class="Profile-player--filter competitive-filter">Competitive Play</blz-tab-control>
```

### Hero Statistics Structure
Each hero has these selectors:
- `data-hero-id="lucio"` - Hero identifier
- `data-progress="100"` - Progress bar percentage
- Hero name in text: "Lúcio"
- Time played: "44:28:48"

### Metrics Available
From dropdown options (data-category-id values):
- `0x0860000000000021` - Time Played
- `0x0860000000000039` - Games Won
- `0x08600000000003D1` - Win Percentage
- `0x08600000000001BB` - Weapon Accuracy - Best in Game
- `0x08600000000003D2` - Eliminations per Life
- `0x0860000000000223` - Kill Streak - Best
- `0x0860000000000346` - Multikill - Best
- `0x08600000000004D4` - Eliminations - Avg per 10 Min
- `0x08600000000004D3` - Deaths - Avg per 10 Min
- `0x08600000000004D5` - Final Blows - Avg per 10 Min
- `0x08600000000004DA` - Solo Kills - Avg per 10 Min
- `0x08600000000004D8` - Objective Kills - Avg per 10 Min
- `0x08600000000004D9` - Objective Time - Avg per 10 Min
- `0x08600000000004BD` - Hero Damage Done - Avg per 10 Min
- `0x08600000000004D6` - Healing Done - Avg per 10 Min

## CSS Selectors for Parsing

### Platform Detection
- Active platform: `.Profile-player--filter.is-active`
- PC platform: `#mouseKeyboardFilter`
- Console platform: `#controllerFilter`

### Game Mode Detection
- Quick Play: `.quickPlay-filter.is-active`
- Competitive: `.competitive-filter.is-active`

### Hero Data Extraction
- Hero containers: `.Profile-progressBar`
- Hero ID: `[data-hero-id]`
- Hero name: `.Profile-progressBar-title`
- Time played: `.Profile-progressBar-description`
- Progress value: `[data-progress]`

### Platform Views
- PC stats: `.mouseKeyboard-view.is-active`
- Console stats: `.controller-view.is-active`
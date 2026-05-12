---
name: KAS Bon
colors:
  primary: "oklch(0.455 0.188 13.697)"
  primary-foreground: "oklch(0.969 0.015 12.422)"
  secondary: "oklch(0.274 0.006 286.033)"
  secondary-foreground: "oklch(0.985 0 0)"
  surface: "oklch(0.141 0.005 285.823)"
  on-surface: "oklch(0.985 0 0)"
  muted: "oklch(0.274 0.006 286.033)"
  muted-foreground: "oklch(0.705 0.015 286.067)"
  accent: "oklch(0.274 0.006 286.033)"
  accent-foreground: "oklch(0.985 0 0)"
  destructive: "oklch(0.704 0.191 22.216)"
  border: "oklch(1 0 0 / 10%)"
  input: "oklch(1 0 0 / 15%)"
  ring: "oklch(0.552 0.016 285.938)"
  card: "oklch(0.21 0.006 285.885)"
  card-foreground: "oklch(0.985 0 0)"
  popover: "oklch(0.21 0.006 285.885)"
  popover-foreground: "oklch(0.985 0 0)"
  sidebar: "oklch(0.21 0.006 285.885)"
  sidebar-foreground: "oklch(0.985 0 0)"
  sidebar-primary: "oklch(0.645 0.246 16.439)"
  chart-1: "oklch(0.871 0.15 154.449)"
  chart-2: "oklch(0.723 0.219 149.579)"
  chart-3: "oklch(0.627 0.194 149.214)"
  chart-4: "oklch(0.527 0.154 150.069)"
  chart-5: "oklch(0.448 0.119 151.328)"
typography:
  heading:
    fontFamily: "JetBrains Mono Variable, monospace"
    fontWeight: 600
  body-md:
    fontFamily: "JetBrains Mono Variable, monospace"
    fontSize: 12px
    fontWeight: 400
  body-lg:
    fontFamily: "JetBrains Mono Variable, monospace"
    fontSize: 14px
    fontWeight: 400
  label:
    fontFamily: "JetBrains Mono Variable, monospace"
    fontSize: 12px
    fontWeight: 500
    lineHeight: 1
  button:
    fontFamily: "JetBrains Mono Variable, monospace"
    fontSize: 12px
    fontWeight: 500
  sidebar-heading:
    fontFamily: "JetBrains Mono Variable, monospace"
    fontSize: 18px
    fontWeight: 600
rounded:
  none: 0px
  sm: 6px
  md: 8px
  lg: 10px
  xl: 14px
  2xl: 18px
  3xl: 22px
  4xl: 26px
---

# Design System

## Overview

A dark-first, monospace interface for an RBAC-based authentication and authorization management system.
Brutalist aesthetic with sharp corners, compact spacing, and a warm red-orange accent against cool zinc neutrals.
High information density, low visual noise, technical precision.

**Default theme: Dark.** Light mode is supported but dark is the primary experience.

## Colors

All colors use the OKLCH color space for perceptual uniformity. The palette is built on a cool zinc (blue-gray hue ~286) neutral base with a warm red-orange primary accent (hue ~14-17).

### Semantic Colors (Dark Mode — Default)

| Token | Value | Usage |
|---|---|---|
| **Primary** | `oklch(0.455 0.188 13.697)` | CTAs, active states, key interactive elements, brand accent |
| **Primary foreground** | `oklch(0.969 0.015 12.422)` | Text/icons on primary backgrounds |
| **Surface** | `oklch(0.141 0.005 285.823)` | Page backgrounds |
| **On-surface** | `oklch(0.985 0 0)` | Primary text on dark surfaces |
| **Card** | `oklch(0.21 0.006 285.885)` | Card and elevated surface backgrounds |
| **Card foreground** | `oklch(0.985 0 0)` | Text on card surfaces |
| **Secondary** | `oklch(0.274 0.006 286.033)` | Supporting UI, chips, secondary actions |
| **Secondary foreground** | `oklch(0.985 0 0)` | Text on secondary backgrounds |
| **Muted** | `oklch(0.274 0.006 286.033)` | Muted/subtle backgrounds |
| **Muted foreground** | `oklch(0.705 0.015 286.067)` | Secondary/disabled text |
| **Accent** | `oklch(0.274 0.006 286.033)` | Hover states, highlighted items |
| **Accent foreground** | `oklch(0.985 0 0)` | Text on accent backgrounds |
| **Destructive** | `oklch(0.704 0.191 22.216)` | Error states, delete actions, validation errors |
| **Border** | `oklch(1 0 0 / 10%)` | Borders and dividers (white at 10% opacity) |
| **Input** | `oklch(1 0 0 / 15%)` | Input field borders (white at 15% opacity) |
| **Ring** | `oklch(0.552 0.016 285.938)` | Focus rings and outlines |

### Light Mode Overrides

| Token | Value | Usage |
|---|---|---|
| **Surface** | `oklch(1 0 0)` | Pure white page backgrounds |
| **On-surface** | `oklch(0.141 0.005 285.823)` | Near-black text |
| **Primary** | `oklch(0.514 0.222 16.935)` | Slightly brighter red-orange |
| **Card** | `oklch(1 0 0)` | White card backgrounds |
| **Secondary** | `oklch(0.967 0.001 286.375)` | Very light gray |
| **Muted** | `oklch(0.967 0.001 286.375)` | Very light gray |
| **Muted foreground** | `oklch(0.552 0.016 285.938)` | Medium gray text |
| **Border** | `oklch(0.92 0.004 286.32)` | Light gray borders |
| **Input** | `oklch(0.92 0.004 286.32)` | Light gray input borders |

### Chart Palette

A green monochromatic palette for data visualization, ranging from light to dark:

| Token | Value | Usage |
|---|---|---|
| **Chart 1** | `oklch(0.871 0.15 154.449)` | Lightest — primary data series |
| **Chart 2** | `oklch(0.723 0.219 149.579)` | Light-medium — second series |
| **Chart 3** | `oklch(0.627 0.194 149.214)` | Medium — third series |
| **Chart 4** | `oklch(0.527 0.154 150.069)` | Medium-dark — fourth series |
| **Chart 5** | `oklch(0.448 0.119 151.328)` | Darkest — fifth series |

### Sidebar Colors

| Token | Dark | Light |
|---|---|---|
| **Sidebar bg** | `oklch(0.21 0.006 285.885)` | `oklch(0.985 0 0)` |
| **Sidebar foreground** | `oklch(0.985 0 0)` | `oklch(0.141 0.005 285.823)` |
| **Sidebar primary** | `oklch(0.645 0.246 16.439)` | `oklch(0.586 0.253 17.585)` |
| **Sidebar primary fg** | `oklch(0.969 0.015 12.422)` | `oklch(0.969 0.015 12.422)` |
| **Sidebar accent** | `oklch(0.274 0.006 286.033)` | `oklch(0.967 0.001 286.375)` |
| **Sidebar border** | `oklch(1 0 0 / 10%)` | `oklch(0.92 0.004 286.32)` |

## Typography

Single font family for all contexts: **JetBrains Mono Variable** (monospace). This gives the interface a technical, developer-oriented feel.

### Scale

| Role | Size | Weight | Line Height | Usage |
|---|---|---|---|---|
| Sidebar heading | 18px | 600 (semi-bold) | Default | App/brand name in sidebar |
| Card title | 14px | 500 (medium) | Default | Card headings |
| Body large | 14px | 400 (regular) | Default | Descriptions, secondary body |
| Body | 12px | 400 (regular) | Default | Primary body text, menu items |
| Label | 12px | 500 (medium) | 1 | Form labels, field labels |
| Button | 12px | 500 (medium) | Default | Button text |
| Tooltip | 12px | 400 (regular) | Default | Tooltip text |

### Font Stack

```css
font-family: 'JetBrains Mono Variable', monospace;
```

Load via:
```html
<link rel="stylesheet" href="https://cdn.jsdelivr.net/npm/@fontsource-variable/jetbrains-mono/index.css">
```

Or install the npm package: `@fontsource-variable/jetbrains-mono`

## Spacing & Sizing

Base unit: **4px**. All spacing uses multiples of 4px.

### Component Sizing

| Component | Height | Padding | Notes |
|---|---|---|---|
| Button default | 32px | 10px horizontal | Gap: 6px between icon and text |
| Button small | 28px | 10px horizontal | Gap: 4px |
| Button extra-small | 24px | 8px horizontal | Gap: 4px |
| Button large | 36px | 10px horizontal | Gap: 6px |
| Icon button | 32x32px | — | Square, centered icon |
| Icon button xs | 24x24px | — | Square, centered icon |
| Icon button sm | 28x28px | — | Square, centered icon |
| Icon button lg | 36x36px | — | Square, centered icon |
| Input | 32px | 10px horizontal, 4px vertical | — |
| Checkbox | 16x16px | — | — |
| Icon (default) | 16px | — | Phosphor icon size |

### Layout Spacing

| Context | Padding |
|---|---|
| Page content | 24px |
| Card body | 16px |
| Card body (small) | 12px |
| Card header/footer | 16px top/bottom |
| Sheet header/footer | 16px |
| Sidebar header/footer | 8px |
| Sidebar group | 8px |
| Tooltip | 12px horizontal, 6px vertical |

### Sidebar Dimensions

| Property | Value |
|---|---|
| Desktop width | 256px (16rem) |
| Mobile width | 288px (18rem) |
| Collapsed (icon-only) width | 48px (3rem) |

## Border Radius

All components use **0px radius (sharp corners)**. This is a deliberate design choice — the brutalist aesthetic forbids rounded corners on any interactive or container element.

The radius scale exists for potential future use:

| Token | Value |
|---|---|
| none | 0px |
| sm | 6px |
| md | 8px |
| lg | 10px |
| xl | 14px |
| 2xl | 18px |
| 3xl | 22px |
| 4xl | 26px |

Base radius variable: `--radius: 0.625rem` (10px). Scale values are calculated as multipliers of this base.

## Borders & Elevation

### Borders

- **Default**: 1px solid using the `border` token color
- **Inputs**: 1px solid using the `input` token color
- **Cards**: No border — uses a 1px ring (`ring-1`) at 10% foreground opacity for visual separation
- **Separators**: 1px height/width using `border` color

### Elevation

No traditional box-shadows. Visual hierarchy is achieved through:

- **Background contrast**: Card surfaces use `card` token, slightly lighter than `background`
- **Ring outlines**: Cards use `ring-1 ring-foreground/10` (1px outline at 10% opacity)
- **Sheet/overlay**: `shadow-lg` only on floating sheet content
- **Sidebar (floating)**: `shadow-sm` on the inner container

## States

### Focus

```
border: [ring color]
outline: 1px solid [ring color at 50% opacity]
```

Focus is indicated by a border color change to `ring` and a 1px ring at 50% opacity. Only appears on keyboard focus (`:focus-visible`), not on click.

### Disabled

```
opacity: 0.5
pointer-events: none
```

### Invalid / Error

```
border: [destructive color]
ring: 1px solid [destructive at 20% opacity]
/* Dark mode */
border: [destructive at 50% opacity]
ring: 1px solid [destructive at 40% opacity]
```

### Active (Press)

Buttons without popups receive a 1px downward translate on press:
```
transform: translateY(1px)
```

## Animation & Transitions

### Duration & Easing

| Context | Duration | Easing |
|---|---|---|
| Sheet overlay | 100ms | default |
| Sheet content | 200ms | ease-in-out |
| Sidebar collapse/expand | 200ms | ease-linear |
| Button | — | transition-all |
| Input | — | transition-colors |

### Animations

| Animation | Usage |
|---|---|
| Fade in/out | Sheet, tooltip open/close (opacity 0 → 1) |
| Slide in/out | Sheet directional entry (±10px slide), tooltip (±2px slide) |
| Zoom in/out | Tooltip scale (95% → 100%) |
| Pulse | Skeleton loading placeholders |
| Spin | Loading spinners |

### Skeleton Loading

```
background: [muted color]
border-radius: 0px
animation: pulse (2s cubic-bezier(0.4, 0, 0.6, 1) infinite)
```

## Components

### Button

6 variants × 8 sizes. Uses a compound component pattern (supports polymorphic rendering).

| Variant | Background | Foreground | Border |
|---|---|---|---|
| Default | Primary | Primary foreground | None |
| Secondary | Secondary | Secondary foreground | None |
| Outline | Transparent | Foreground | Border color |
| Ghost | Transparent | Foreground | None |
| Destructive | Destructive | Primary foreground | None |
| Link | Transparent | Primary | None (underline on hover) |

### Input

- Height: 32px
- Border: 1px `input` color
- Background: transparent (slightly tinted in dark mode: `input` at 30% opacity)
- Supports file input variant
- Focus: border + ring outline
- Invalid: destructive border + ring

### Card

7 sub-components: Card, CardHeader, CardTitle, CardDescription, CardAction, CardContent, CardFooter.

- No border, no shadow — uses `ring-1 ring-foreground/10` for visual edge
- Supports `size="sm"` variant with tighter padding
- Uses container queries for responsive internal layout

### Sheet (Drawer/Modal)

Radix Dialog-based sliding panel from any edge (top, right, bottom, left).

- Overlay: background at 80% opacity + `backdrop-blur-xs`
- Content: card background, rounded-none
- Border on the leading edge (side-dependent)
- Optional close button (enabled by default)

### Sidebar

Full navigation sidebar with 3 layout variants:

| Variant | Description |
|---|---|
| Default | Flush to viewport edge |
| Floating | Offset with spacing and `shadow-sm` |
| Inset | Embedded within content area |

Collapsible modes:
- **Offcanvas**: Slides in/out (default)
- **Icon**: Collapses to icon-only (48px wide)
- **None**: Always expanded

Mobile behavior: Replaces with a Sheet (drawer) below 768px.

Keyboard shortcut: `Ctrl+B` to toggle.

### Checkbox

- Size: 16x16px
- Border: 1px `input` color
- Checked: primary background, primary-foreground check icon
- Focus: ring outline
- Transition: `transition-colors`

### Tooltip

- Inverted colors: foreground background, background text
- Arrow: 6px rotated square (45deg)
- Zero delay (appears immediately on hover)
- Content padding: 12px horizontal, 6px vertical

### Toast (Sonner)

Theme-aware toast notifications with custom icons per state (success, error, info, warning).
Integrates with theme system for automatic dark/light mode colors.

### Separator

- 1px height (horizontal) or 1px width (vertical)
- Background: `border` color
- Supports both orientations

## Iconography

**Icon library: Phosphor Icons** (`@phosphor-icons/react`)

- Default size: 16px (`size-4`)
- Weight: Regular (default)
- All icons should use Phosphor for consistency
- Do not mix icon libraries

Vanilla HTML alternative: use Phosphor Icons web component or SVG sprites.

## Responsive Breakpoints

| Breakpoint | Width | Usage |
|---|---|---|
| sm | 640px | Sheet max-width constraints |
| md | 768px | Sidebar visibility (mobile ↔ desktop), mobile detection threshold |

No custom breakpoints — use standard CSS media queries:

```css
@media (min-width: 640px) { /* sm */ }
@media (min-width: 768px) { /* md */ }
```

## CSS Custom Properties Reference

Define these as CSS custom properties on `:root` for light mode and `.dark` class for dark mode:

```css
:root {
  --background: oklch(1 0 0);
  --foreground: oklch(0.141 0.005 285.823);
  --primary: oklch(0.514 0.222 16.935);
  --primary-foreground: oklch(0.969 0.015 12.422);
  --secondary: oklch(0.967 0.001 286.375);
  --secondary-foreground: oklch(0.21 0.006 285.885);
  --muted: oklch(0.967 0.001 286.375);
  --muted-foreground: oklch(0.552 0.016 285.938);
  --accent: oklch(0.967 0.001 286.375);
  --accent-foreground: oklch(0.21 0.006 285.885);
  --destructive: oklch(0.577 0.245 27.325);
  --border: oklch(0.92 0.004 286.32);
  --input: oklch(0.92 0.004 286.32);
  --ring: oklch(0.705 0.015 286.067);
  --card: oklch(1 0 0);
  --card-foreground: oklch(0.141 0.005 285.823);
  --popover: oklch(1 0 0);
  --popover-foreground: oklch(0.141 0.005 285.823);
  --radius: 0.625rem;
  /* Chart */
  --chart-1: oklch(0.871 0.15 154.449);
  --chart-2: oklch(0.723 0.219 149.579);
  --chart-3: oklch(0.627 0.194 149.214);
  --chart-4: oklch(0.527 0.154 150.069);
  --chart-5: oklch(0.448 0.119 151.328);
  /* Sidebar */
  --sidebar: oklch(0.985 0 0);
  --sidebar-foreground: oklch(0.141 0.005 285.823);
  --sidebar-primary: oklch(0.586 0.253 17.585);
  --sidebar-primary-foreground: oklch(0.969 0.015 12.422);
  --sidebar-accent: oklch(0.967 0.001 286.375);
  --sidebar-accent-foreground: oklch(0.21 0.006 285.885);
  --sidebar-border: oklch(0.92 0.004 286.32);
  --sidebar-ring: oklch(0.705 0.015 286.067);
}

.dark {
  --background: oklch(0.141 0.005 285.823);
  --foreground: oklch(0.985 0 0);
  --card: oklch(0.21 0.006 285.885);
  --card-foreground: oklch(0.985 0 0);
  --popover: oklch(0.21 0.006 285.885);
  --popover-foreground: oklch(0.985 0 0);
  --primary: oklch(0.455 0.188 13.697);
  --primary-foreground: oklch(0.969 0.015 12.422);
  --secondary: oklch(0.274 0.006 286.033);
  --secondary-foreground: oklch(0.985 0 0);
  --muted: oklch(0.274 0.006 286.033);
  --muted-foreground: oklch(0.705 0.015 286.067);
  --accent: oklch(0.274 0.006 286.033);
  --accent-foreground: oklch(0.985 0 0);
  --destructive: oklch(0.704 0.191 22.216);
  --border: oklch(1 0 0 / 10%);
  --input: oklch(1 0 0 / 15%);
  --ring: oklch(0.552 0.016 285.938);
  --sidebar: oklch(0.21 0.006 285.885);
  --sidebar-foreground: oklch(0.985 0 0);
  --sidebar-primary: oklch(0.645 0.246 16.439);
  --sidebar-primary-foreground: oklch(0.969 0.015 12.422);
  --sidebar-accent: oklch(0.274 0.006 286.033);
  --sidebar-accent-foreground: oklch(0.985 0 0);
  --sidebar-border: oklch(1 0 0 / 10%);
  --sidebar-ring: oklch(0.552 0.016 285.938);
}
```

## Do's and Don'ts

- Do use the primary red-orange color sparingly — only for the most important action per view
- Do maintain sharp corners (0px radius) on all components for visual consistency
- Do use `ring-1` outlines instead of box-shadows for card and surface boundaries
- Do use OKLCH color values to ensure perceptually uniform color relationships
- Do apply the monospace font (JetBrains Mono) consistently across all text contexts
- Do use the zinc-based neutral palette for all non-accent UI elements
- Do keep focus indicators visible — use the ring color at 50% opacity for focus states
- Do use Phosphor icons exclusively at 16px default size
- Don't mix rounded and sharp corners in the same view
- Don't use traditional box-shadows for elevation — rely on background contrast and ring outlines
- Don't introduce additional font families — JetBrains Mono is the sole typeface
- Don't use the destructive color for anything other than errors and irreversible actions
- Don't apply border-radius to buttons, inputs, cards, or any interactive element
- Don't add color beyond the defined palette — extend the system tokens, don't hardcode hex values
- Don't use opacity below 50% for disabled states — maintain minimum readability

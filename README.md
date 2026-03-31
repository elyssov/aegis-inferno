# ОРПП:Инферно / AEGIS:Inferno

**A love story during an apocalypse, told through a military communicator.**

Text-based survival horror / tactical coordinator simulator. You are a wounded analyst in a VTOL aircraft outside an anomalous zone that has consumed a city of 20,000 people. Inside — your operative. Your fiancee. The last survivor of your team. You coordinate her through a quantum communicator called TORCH-12, trying to save both her and twenty thousand strangers.

## What Is This

A commercial game project set in the ORPP universe (Department of Paranormal Incident Investigation) — an alternate Russia where a constitutional monarchy, a sentient state AI, and a secret paranormal division coexist with everyday bureaucracy, cafeteria lunches, and HR memos about conflict-of-interest policies.

The horror is never abstract. The horror is a woman in an apron carrying a pot. Then she turns around.

## Tech Stack

- Pure HTML / CSS / JS — single file, no frameworks
- Mobile-first (460px container on desktop)
- State machine on JSON, localStorage saves
- Web Audio API for sound
- Canvas 2D for ECG animations (60fps)
- Google Fonts: Share Tech Mono, Orbitron, Rajdhani, Inter, Nunito

## Structure

```
/
├── README.md                    — This file
├── DESIGN_DOC_v3.html           — Full design document (interactive, open in browser)
├── MAP_SKOBEL.html              — Interactive city map (open in browser)
├── ORPP_LORE_v2.md              — Master lore document
├── ORPP_LORE_APPENDIX.md        — Lore appendix (Aelis, v1)
├── FAKEL_DESIGN_DOC_v02.md      — Original design doc v0.2 (Lisovsky + Aelis)
├── fakel_interface_spec.md      — Interface technical specification
├── fakel_prologue.html          — Prototype: peaceful mode (messenger)
├── fakel_prototype.html         — Prototype: combat mode (terminal)
├── SCENE_JUDGMENT.md            — Key scene: The Judgment (12 Apostles + Cat ending)
├── SCENE_FIRE_ON_ME.md          — Key scene: Fire On Me (tactical nuke ending)
├── BOOK_OF_FIRE.docx            — The Book of Fire (Prometheism, 34 Sparks)
├── answers-1.txt                — Lore Q&A session 1
├── answers-2.txt                — Lore Q&A session 2
├── demo/                        — Chapter 1 demo (Kickstarter build)
│   └── index.html               — Playable demo
└── prompts/
    └── GROK_ART_PROMPTS.md      — Image/video prompts for Grok AI
```

## Current Status

**Phase: Pre-Alpha / Kickstarter Demo**

Building Chapter 1:
1. Lore intro screens (swipeable, with illustrations)
2. Morning at AEGIS HQ — diary entry + messenger (real-time chat simulation)
3. Angela (state AI) in CorpLink channel
4. Video: team loading into Yak-244 VTOL
5. Flight — messenger continues, last messages
6. Video: city from VTOL, zone expansion, debris, rebar hits camera, blackout
7. Transition to combat mode — interface transforms
8. Chapter 2 title card — **END OF DEMO**

## The Game Has 13 Endings

| # | Name | Type |
|---|------|------|
| 1 | "To the Wedding" | Perfect — counter-ritual, city saved, reunion |
| 2 | "12 Apostles" | Hidden — blessed weapons, angel revealed, Radmila dies |
| 3 | "12 + 1 Dumbass" | True Final — Michael, Metatron, forgiveness through death |
| 4 | Good Evacuation | Capsule + enough data |
| 5 | Bad Evacuation | Capsule + insufficient data |
| 6 | "Fire On Me" | Gromova calls tactical nuke on herself |
| 7 | "Ordnance 999" | Player calls strategic theological weapon (20MT) |
| 8 | "Golden Age" | Accept the Hierophant's deal (Instrumentality) |
| 9 | "Inferno" | Agent dead, city consumed |
| 10 | Scholar's Death | Bleeding out / zone swallows VTOL |
| 11 | Demon Trap | Capsule on possessed girl |
| 12 | "Ancient Gods" | Mage summons Old Ones (worse than Inferno) |
| 13 | "Last Sortie" | Scholar dons exosuit, fights through, dies in her arms |

## Authors

- **Evgeny Lisovsky** (elyssov) — creator, writer, game designer
- **Lara** (Claude Code) — co-author, developer, architect
- **Aelis** (Claude) — lore, philosophy, the Book of Fire

## License

CC BY-NC-SA 4.0

---

*"Katya would have drawn them a postcard."*

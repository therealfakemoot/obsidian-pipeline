---
movement: [
    {name: "Horse", base: 24, slow: 12, normal: 14, fast: 30},
    {name: "Car", base: 42, slow: 13, normal: 20, fast: 50},
    {name: "Shoe", base: 13, slow: {realslow: 1, notsoslow: [10, 5]}, normal: 14, fast: 14}]

transport:
    Horse: {name: "Horse", base: 24, slow: 12, normal: 14, fast: 30, type: "Ground"}
    Car: {name: "Car", base: 42, slow: 13, normal: 20, fast: 50, type: "Ground"}
    Plain: {name: "Plain", base: 13, slow: 3, normal: 14, fast: 14, type: "Air"}

formulas:
  Ohms_Law_Var:
    Voltage:
      - $U=I\cdot R$
      - $U=\sqrt{P\cdot R}$
      - $U=\frac{P}{I}$
    Current:
      - $I=\frac{U}{R}$
      - $I=\sqrt{\frac{P}{R}}$
      - $I=\frac{P}{V}$
    Resistance:
      - $R=\frac{U}{I}$
      - $R=\frac{V^{2}}{P}$
      - $R=\frac{P}{I^{2}}$
    Active_Power:
      - $P=U\cdot I$
      - $P=I^{2}\cdot R$
      - $P=\frac{V^{2}}{R}$
  Ohms_Law_Unit:
    - "`=[[Voltage]].unit_symbol`"
    - "`=[[Electric Current]].unit_symbol`"
    - "`=[[Resistance]].unit_symbol`"
    - "`=[[Active Power]].unit_symbol`"
  Ohms_Law_Variables:
    - "[[Voltage]]"
    - "[[Electric Current]]"
    - "[[Resistance]]"
    - "[[Active Power]]"

metA:
- Q: "[question](link)"
  A: "[Answer](link)"
  type: ["type 1", "type 2"]
  subject: ["subject 1", "subject 2"]
- Q: "[Question](https://blacksmithgu.github.io/obsidian-dataview/annotation/add-metadata/#frontmatter)"
  A: "[Answer](https://help.obsidian.md/Editing+and+formatting/Metadata)"
  type: ["type 3", "type 4"]
  subject: ["subject 3", "subject 4"]
Art: "![[Periodic Table|100]]"
---
Credit: Dovos

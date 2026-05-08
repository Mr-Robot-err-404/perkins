package theme

var Kanagawa = buildPalette("Kanagawa", paletteData{
	fg: []string{
		// Page 1: warm — samuraiRed, waveRed, surimiOrange, roninYellow, autumnYellow, carpYellow, fujiWhite, boatYellow2
		"#E82424", "#E46876", "#FFA066", "#FF9E3B",
		"#DCA561", "#E6C384", "#DCD7BA", "#C0A36E",

		// Page 2: cool — springGreen, autumnGreen, waveAqua1, waveAqua2, lightBlue, springBlue, crystalBlue, springViolet2
		"#98BB6C", "#76946A", "#6A9589", "#7AA89F",
		"#A3D4D5", "#7FB4CA", "#7E9CD8", "#9CABCA",

		// Page 3: oniViolet, springViolet1, sakuraPink, peachRed — sumiInk6, katanaGray, sumiInk5, winterBlue
		"#957FB8", "#938AA9", "#D27E99", "#FF5D62",
		"#54546D", "#717C7C", "#363646", "#252535",
	},
	bg: []string{
		// Page 1: bright/mid hue sweep — carpYellow, surimiOrange, sakuraPink, lotusPink(bg), springGreen, waveAqua2, crystalBlue, oniViolet
		"#E6C384", "#FFA066", "#D27E99", "#D798A6",
		"#98BB6C", "#7AA89F", "#7E9CD8", "#957FB8",

		// Page 2: muted/dark hue sweep — autumnRed, autumnYellow(dark), autumnGreen, waveAqua1, dragonBlue, waveBlue2, katanaGray, sumiInk5
		"#C34043", "#49443C", "#76946A", "#6A9589",
		"#658594", "#3B4C77", "#717C7C", "#363646",

		// Page 3: ink darks — waveBlue1, winterGreen, winterRed, winterBlue, sumiInk3, sumiInk2, sumiInk1, sumiInk0
		"#223249", "#2B3328", "#43242B", "#252535",
		"#1F1F28", "#161628", "#16161D", "#0D0C0C",
	},
})

package theme

var TokyoNight = buildPalette("TokyoNight", paletteData{
	fg: []string{
		"#FF9E64", "#E5C890", "#E0AF68", "#D19A66",
		"#C3E88D", "#A6D189", "#9ECE6A", "#73DAC8",

		"#89DDFF", "#7DCFFF", "#66D9EF", "#C0CAF5",
		"#89B4FA", "#82AAFF", "#A9B1D6", "#7AA2F7",

		"#565F89", "#414868", "#343B58", "#2A2E49",
		"#BB9AF7", "#FCA7EA", "#FF757F", "#F7768E",
	},
	bg: []string{
		"#C3E88D", "#89DDFF", "#C0CAF5", "#66D9EF",
		"#73DAC8", "#7DCFFF", "#9ECE6A", "#E0AF68",

		"#FF9E64", "#A9B1D6", "#82AAFF", "#BB9AF7",
		"#D19A66", "#7AA2F7", "#FF757F", "#F7768E",

		"#565F89", "#414868", "#343B58", "#292E42",
		"#24283B", "#1F2335", "#1A1B26", "#0D0E17",
	},
})

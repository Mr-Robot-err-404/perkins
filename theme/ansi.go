package theme

var ANSI = buildPalette("ANSI 16-bit", paletteData{
	fg: []string{
		"#AA0000", "#00AA00", "#AA5500", "#0000AA",
		"#AA00AA", "#00AAAA", "#AAAAAA", "#555555",

		"#FF5555", "#55FF55", "#FFFF55", "#5555FF",
		"#FF55FF", "#55FFFF", "#FFFFFF", "#000000",
	},
	bg: []string{
		"#AA0000", "#00AA00", "#AA5500", "#0000AA",
		"#AA00AA", "#00AAAA", "#AAAAAA", "#555555",

		"#FF5555", "#55FF55", "#FFFF55", "#5555FF",
		"#FF55FF", "#55FFFF", "#FFFFFF", "#000000",
	},
})

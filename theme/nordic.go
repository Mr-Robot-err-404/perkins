package theme

var Nordic = Palette{
	Name: "Nordic",
	Foreground: []Color{
		// page 1 - oranges, yellows, greens, teals
		{Ansi: "38;2;220;165;120", Display: "#DCA578"},
		{Ansi: "38;2;208;135;112", Display: "#D08770"},
		{Ansi: "38;2;245;213;155", Display: "#F5D59B"},
		{Ansi: "38;2;235;203;139", Display: "#EBCB8B"},
		{Ansi: "38;2;178;205;155", Display: "#B2CD9B"},
		{Ansi: "38;2;163;190;140", Display: "#A3BE8C"},
		{Ansi: "38;2;171;216;214", Display: "#ABD8D6"},
		{Ansi: "38;2;136;192;208", Display: "#88C0D0"},
		// page 2 - blues, whites
		{Ansi: "38;2;143;188;187", Display: "#8FBCBB"},
		{Ansi: "38;2;150;200;220", Display: "#96C8DC"},
		{Ansi: "38;2;236;239;244", Display: "#ECEFF4"},
		{Ansi: "38;2;229;233;240", Display: "#E5E9F0"},
		{Ansi: "38;2;216;222;233", Display: "#D8DEE9"},
		{Ansi: "38;2;129;161;193", Display: "#81A1C1"},
		{Ansi: "38;2;110;170;198", Display: "#6EAAC6"},
		{Ansi: "38;2;94;129;172",  Display: "#5E81AC"},
		// page 3 - dark blues, purples, reds
		{Ansi: "38;2;76;86;106",   Display: "#4C566A"},
		{Ansi: "38;2;67;76;94",    Display: "#434C5E"},
		{Ansi: "38;2;59;66;82",    Display: "#3B4252"},
		{Ansi: "38;2;46;52;64",    Display: "#2E3440"},
		{Ansi: "38;2;196;158;188", Display: "#C49EBC"},
		{Ansi: "38;2;180;142;173", Display: "#B48EAD"},
		{Ansi: "38;2;210;115;125", Display: "#D2737D"},
		{Ansi: "38;2;191;97;106",  Display: "#BF616A"},
	},
	Background: []Color{
		// page 1 - light
		{Ansi: "48;2;229;233;240", Display: "#E5E9F0"},
		{Ansi: "48;2;216;222;233", Display: "#D8DEE9"},
		{Ansi: "48;2;245;213;155", Display: "#F5D59B"},
		{Ansi: "48;2;171;216;214", Display: "#ABD8D6"},
		{Ansi: "48;2;235;203;139", Display: "#EBCB8B"},
		{Ansi: "48;2;178;205;155", Display: "#B2CD9B"},
		{Ansi: "48;2;150;200;220", Display: "#96C8DC"},
		{Ansi: "48;2;136;192;208", Display: "#88C0D0"},
		// page 2 - mid
		{Ansi: "48;2;163;190;140", Display: "#A3BE8C"},
		{Ansi: "48;2;143;188;187", Display: "#8FBCBB"},
		{Ansi: "48;2;220;165;120", Display: "#DCA578"},
		{Ansi: "48;2;196;158;188", Display: "#C49EBC"},
		{Ansi: "48;2;110;170;198", Display: "#6EAAC6"},
		{Ansi: "48;2;129;161;193", Display: "#81A1C1"},
		{Ansi: "48;2;180;142;173", Display: "#B48EAD"},
		{Ansi: "48;2;208;135;112", Display: "#D08770"},
		// page 3 - dark
		{Ansi: "48;2;94;129;172",  Display: "#5E81AC"},
		{Ansi: "48;2;191;97;106",  Display: "#BF616A"},
		{Ansi: "48;2;94;106;130",  Display: "#5E6A82"},
		{Ansi: "48;2;76;86;106",   Display: "#4C566A"},
		{Ansi: "48;2;67;76;94",    Display: "#434C5E"},
		{Ansi: "48;2;59;66;82",    Display: "#3B4252"},
		{Ansi: "48;2;46;52;64",    Display: "#2E3440"},
		{Ansi: "48;2;36;41;51",    Display: "#242933"},
	},
}

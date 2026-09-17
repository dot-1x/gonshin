package style

import "testing"

func TestCN(t *testing.T) {
	tests := []struct {
		name string
		in   []string
		want string
	}{
		{"padding override", []string{"py-6 py-0"}, "py-0"},
		{"padding x override", []string{"px-6 px-0"}, "px-0"},
		{"padding all over x", []string{"px-6 p-0"}, "p-0"},
		{"padding x after all stays", []string{"p-0 px-6"}, "p-0 px-6"},
		{"padding axes coexist", []string{"px-4 py-3"}, "px-4 py-3"},
		{"gap override", []string{"gap-6 gap-0"}, "gap-0"},
		{"size does not override w/h", []string{"size-8 h-32 w-32"}, "size-8 h-32 w-32"},
		{"text size override", []string{"text-xs text-[0.65rem]"}, "text-[0.65rem]"},
		{"text colour keeps size", []string{"text-xs text-muted-foreground"}, "text-xs text-muted-foreground"},
		{"text size keeps colour", []string{"text-muted-foreground text-xs"}, "text-muted-foreground text-xs"},
		{"max width override", []string{"sm:max-w-lg sm:max-w-2xl"}, "sm:max-w-2xl"},
		{"bg colour override", []string{"bg-card bg-popover"}, "bg-popover"},
		{"border colour override", []string{"border-transparent border-secondary/20"}, "border-secondary/20"},
		{"border width and colour", []string{"border border-2 border-transparent border-primary"}, "border-2 border-primary"},
		{"bare border side width", []string{"border-t border-border/40"}, "border-t border-border/40"},
		{"rounded override", []string{"rounded-full rounded-md"}, "rounded-md"},
		{"font weight override", []string{"font-medium font-normal"}, "font-normal"},
		{"font weight wins", []string{"font-medium font-normal font-medium"}, "font-medium"},
		{"font family separate", []string{"font-medium font-mono"}, "font-medium font-mono"},
		{"width override", []string{"w-fit w-64"}, "w-64"},
		{"card merge", []string{
			"bg-card text-card-foreground flex flex-col gap-6 rounded-xl border py-6 shadow-sm",
			"group cursor-pointer overflow-hidden border-border/50 bg-card py-0 backdrop-blur-sm",
		}, "text-card-foreground flex flex-col gap-6 rounded-xl border shadow-sm group cursor-pointer overflow-hidden border-border/50 bg-card py-0 backdrop-blur-sm"},
		{"badge merge", []string{
			"inline-flex items-center justify-center rounded-full border border-transparent px-2 py-0.5 text-xs font-medium",
			"font-mono text-[0.6rem]",
		}, "inline-flex items-center justify-center rounded-full border border-transparent px-2 py-0.5 font-medium font-mono text-[0.6rem]"},
		{"avatar size kept", []string{
			"group/avatar relative flex size-8 shrink-0 overflow-hidden rounded-full select-none data-[size=lg]:size-10 data-[size=sm]:size-6",
			"h-32 w-32 border-2 border-primary",
		}, "group/avatar relative flex size-8 shrink-0 overflow-hidden rounded-full select-none data-[size=lg]:size-10 data-[size=sm]:size-6 h-32 w-32 border-2 border-primary"},
		{"dialog merge", []string{
			"bg-background data-[state=open]:animate-in fixed top-[50%] left-[50%] grid w-full max-w-[calc(100%-2rem)] gap-4 rounded-lg border p-6 shadow-lg duration-200 sm:max-w-lg",
			"max-h-[85vh] gap-0 overflow-hidden border-border/50 bg-popover p-0 sm:max-w-2xl",
		}, "data-[state=open]:animate-in fixed top-[50%] left-[50%] grid w-full max-w-[calc(100%-2rem)] rounded-lg border shadow-lg duration-200 max-h-[85vh] gap-0 overflow-hidden border-border/50 bg-popover p-0 sm:max-w-2xl"},
		{"variant isolation", []string{"text-muted-foreground hover:text-foreground"}, "text-muted-foreground hover:text-foreground"},
		{"dedupe identical", []string{"bg-card bg-card"}, "bg-card"},
		{"empty input", []string{"", "  "}, ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CN(tt.in...); got != tt.want {
				t.Errorf("CN(%q)\n got: %q\nwant: %q", tt.in, got, tt.want)
			}
		})
	}
}

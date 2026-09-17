// Package ui contains templ ports of the shadcn/ui primitives used by the
// Genshin Impact page, with identical Tailwind class strings.
package ui

import "github.com/dotcchix/gonshin/internal/style"

const (
	cardBase        = "bg-card text-card-foreground flex flex-col gap-6 rounded-xl border py-6 shadow-sm"
	cardHeaderBase  = "@container/card-header grid auto-rows-min grid-rows-[auto_auto] items-start gap-2 px-6 has-data-[slot=card-action]:grid-cols-[1fr_auto] [.border-b]:pb-6"
	cardTitleBase   = "leading-none font-semibold"
	cardActionBase  = "col-start-2 row-span-2 row-start-1 self-start justify-self-end"
	cardContentBase = "px-6"
	cardDescBase    = "text-muted-foreground text-sm"
	cardFooterBase  = "flex items-center px-6 [.border-t]:pt-6"

	badgeBase = "inline-flex items-center justify-center rounded-full border border-transparent px-2 py-0.5 text-xs font-medium w-fit whitespace-nowrap shrink-0 [&>svg]:size-3 gap-1 [&>svg]:pointer-events-none focus-visible:border-ring focus-visible:ring-ring/50 focus-visible:ring-[3px] aria-invalid:ring-destructive/20 dark:aria-invalid:ring-destructive/40 aria-invalid:border-destructive transition-[color,box-shadow] overflow-hidden"

	avatarBase         = "group/avatar relative flex size-8 shrink-0 overflow-hidden rounded-full select-none data-[size=lg]:size-10 data-[size=sm]:size-6"
	avatarImageBase    = "aspect-square size-full"
	avatarFallbackBase = "bg-muted text-muted-foreground flex size-full items-center justify-center rounded-full text-sm"

	separatorBase = "bg-border shrink-0 data-[orientation=horizontal]:h-px data-[orientation=horizontal]:w-full data-[orientation=vertical]:h-full data-[orientation=vertical]:w-px"

	tooltipContentBase = "bg-foreground text-background animate-in fade-in-0 zoom-in-95 data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=closed]:zoom-out-95 data-[side=bottom]:slide-in-from-top-2 data-[side=left]:slide-in-from-right-2 data-[side=right]:slide-in-from-left-2 data-[side=top]:slide-in-from-bottom-2 z-50 w-fit origin-(--radix-tooltip-content-transform-origin) rounded-md px-3 py-1.5 text-xs text-balance"

	dialogContentBase = "bg-background data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 data-[state=closed]:zoom-out-95 data-[state=open]:zoom-in-95 fixed top-[50%] left-[50%] z-50 grid w-full max-w-[calc(100%-2rem)] translate-x-[-50%] translate-y-[-50%] gap-4 rounded-lg border p-6 shadow-lg duration-200 sm:max-w-lg"
	dialogOverlayBase = "data-[state=open]:animate-in data-[state=closed]:animate-out data-[state=closed]:fade-out-0 data-[state=open]:fade-in-0 fixed inset-0 z-50 bg-black/50"
)

var badgeVariants = map[string]string{
	"default":     "bg-primary text-primary-foreground [a&]:hover:bg-primary/90",
	"secondary":   "bg-secondary text-secondary-foreground [a&]:hover:bg-secondary/90",
	"destructive": "bg-destructive text-white [a&]:hover:bg-destructive/90 focus-visible:ring-destructive/20 dark:focus-visible:ring-destructive/40 dark:bg-destructive/60",
	"outline":     "border-border text-foreground [a&]:hover:bg-accent [a&]:hover:text-accent-foreground",
	"ghost":       "[a&]:hover:bg-accent [a&]:hover:text-accent-foreground",
	"link":        "text-primary underline-offset-4 [a&]:hover:underline",
}

func badgeClass(variant, class string) string {
	return style.CN(badgeBase, badgeVariants[variant], class)
}

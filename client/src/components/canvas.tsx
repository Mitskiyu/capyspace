"use client";

import StickyNote from "@/components/widgets/sticky-note";
import { Space } from "@/lib/types";

export default function Canvas({ space }: { space: Space }) {
	return (
		<div className="relative h-screen w-full overflow-hidden">
			{space.widgets.map((widget) => {
				switch (widget.type) {
					case "sticky-note":
						return <StickyNote key={widget.id} widget={widget} />;
					default:
						return null;
				}
			})}
		</div>
	);
}

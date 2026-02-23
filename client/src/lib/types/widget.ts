import { z } from "zod";

export const widgetSchema = z.discriminatedUnion("type", [
	z.object({
		id: z.uuid(),
		type: z.literal("sticky-note"),
		x_pos: z.int32(),
		y_pos: z.int32(),
		minimized: z.boolean(),
		data: z.object({
			text: z.string(),
		}),
	}),
]);

export type Widget = z.infer<typeof widgetSchema>;

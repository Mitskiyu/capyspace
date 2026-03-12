import { z } from "zod";
import { widgetSchema } from "./widget";

export const spaceSchema = z.object({
	id: z.uuid(),
	is_private: z.boolean(),
	widgets: z.array(widgetSchema),
});

export type Space = z.infer<typeof spaceSchema>;

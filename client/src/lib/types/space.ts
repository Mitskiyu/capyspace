import { z } from "zod";

export const spaceSchema = z.object({
	id: z.uuid(),
	is_private: z.boolean(),
});

export type Space = z.infer<typeof spaceSchema>;

import { getSpace } from "@/lib/space";

export default async function Page({
	params,
}: {
	params: { username: string };
}) {
	const { username } = await params;
	const result = await getSpace(username);

	switch (result.status) {
		case "not_found":
			return (
				<div className="text-vibrant-cloud text-xl">Space does not exist</div>
			);

		case "private":
			return (
				<div className="text-vibrant-dream text-xl">This space is private</div>
			);

		case "error":
			return <div className="text-vibrant-coral text-xl">{result.message}</div>;

		case "success":
			const { space } = result;
			return (
				<div>
					<div className="text-vibrant-cloud text-xl">
						{username}&apos;s space: {space.id}
					</div>
				</div>
			);

		default:
			return <div></div>;
	}
}

import Image from "next/image";
import { google, logo } from "@/assets";

function AuthForm() {
	return (
		<div className="bg-neutrals-base text-neutrals-text h-86 w-full max-w-78 rounded-4xl sm:h-92">
			<div className="flex flex-col items-center">
				<div className="flex flex-col items-center">
					<Image
						src={logo}
						alt="logo"
						width={56}
						height={56}
						className="mt-4"
					/>
					<h1 className="pt-2 text-xl sm:text-2xl">Welcome to Capyspace</h1>
					<h2 className="text-neutrals-subtext0 text-lg sm:text-xl">
						Sign up or log in
					</h2>
				</div>
				<button className="bg-neutrals-surface0 mt-4 flex h-10 w-11/12 transform flex-row items-center justify-center space-x-2 rounded-2xl p-2 transition duration-200 hover:-translate-y-0.5 hover:cursor-pointer hover:opacity-90 active:opacity-90 sm:h-12">
					<Image src={google} alt="google" width={24} height={24} />
					<span className="text-lg">Continue with Google</span>
				</button>
				<div className="relative mt-4 flex w-11/12 items-center">
					<hr className="bg-vibrant-lilac from-vibrant-lilac via-vibrant-orchid to-neutrals-lilac h-0.5 w-full rounded-full border-none bg-gradient-to-r" />
					<span className="bg-neutrals-base text-vibrant-orchid absolute left-1/2 -translate-x-1/2 px-2 text-base font-medium">
						OR
					</span>
				</div>
				<div className="bg-neutrals-surface0 mt-6 flex h-10 w-11/12 rounded-2xl sm:h-12">
					<input
						placeholder="Enter your email"
						className="w-full rounded-2xl bg-transparent px-4 py-2 text-base outline-none placeholder:text-base sm:text-lg sm:placeholder:text-lg"
					/>
				</div>
				<button className="text-neutrals-text bg-neutrals-surface0 hover:text-vibrant-orchid mt-2 flex h-10 w-11/12 flex-row items-center justify-center space-x-2 rounded-2xl transition duration-200 hover:cursor-pointer hover:opacity-90 sm:h-12">
					<span className="text-base sm:text-lg">Continue</span>
				</button>
			</div>
		</div>
	);
}

export default AuthForm;

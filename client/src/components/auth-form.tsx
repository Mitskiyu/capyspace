"use client";
import { useForm } from "react-hook-form";
import { z } from "zod";
import Image from "next/image";
import { google, logo } from "@/assets";
import { checkEmail } from "@/lib/auth/check-email";
import { zodResolver } from "@hookform/resolvers/zod";
import { useState } from "react";

const schema = z.object({
	email: z.email("Email is invalid"),
});

type Inputs = z.infer<typeof schema>;
type State = "email" | "signup" | "login";

function AuthForm() {
	const [state, setState] = useState<State>("email");

	const {
		register,
		handleSubmit,
		clearErrors,
		formState: { errors, isValid },
	} = useForm<Inputs>({
		resolver: zodResolver(schema),
		mode: "onTouched",
	});

	async function onSubmit(data: Inputs) {
		const exists = await checkEmail(data.email);
		setState(exists ? "login" : "signup");
	}

	return (
		<div className="bg-neutrals-base text-neutrals-text min-h-83 w-full max-w-78 rounded-4xl sm:min-h-90">
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
				<button className="bg-neutrals-surface0 mt-4 flex h-10 w-11/12 transform cursor-pointer flex-row items-center justify-center space-x-2 rounded-2xl p-2 transition duration-200 hover:-translate-y-0.5 hover:opacity-90 active:opacity-90 sm:h-12">
					<Image src={google} alt="google" width={24} height={24} />
					<span className="text-lg">Continue with Google</span>
				</button>
				<div className="relative mt-4 flex w-11/12 items-center">
					<hr className="bg-vibrant-lilac from-vibrant-lilac via-vibrant-orchid to-neutrals-lilac h-0.5 w-full rounded-full border-none bg-gradient-to-r" />
					<span className="bg-neutrals-base text-vibrant-orchid absolute left-1/2 -translate-x-1/2 px-2 text-base font-medium">
						OR
					</span>
				</div>
				<form
					onSubmit={handleSubmit(onSubmit)}
					className="flex w-full flex-col items-center"
				>
					{errors.email && (
						<div className="text-vibrant-coral mt-2 text-sm">
							{errors.email && <span>{errors.email.message}</span>}
						</div>
					)}
					<input
						{...register("email")}
						onChange={(e) => {
							if (state !== "email") setState("email");
							register("email").onChange(e);
							clearErrors("email");
						}}
						type="text"
						placeholder="Enter your email"
						className={`bg-neutrals-surface0 focus:border-vibrant-bloom border-neutrals-overlay0 h-10 w-11/12 rounded-2xl border-1 px-4 py-2 text-base transition-colors duration-300 ease-out placeholder:text-base focus:border-2 focus:outline-none sm:h-12 sm:text-lg sm:placeholder:text-lg ${errors.email ? "border-vibrant-coral mt-2 border-2" : "mt-4"}`}
					/>

					{state === "login" && <p>Login</p>}

					{state === "signup" && <p>Signup</p>}
					<button
						disabled={!isValid}
						className={`border-neutrals-overlay0 mt-2 mb-4 flex h-10 w-11/12 items-center justify-center rounded-2xl border-1 text-base transition duration-200 sm:h-12 sm:text-lg ${
							isValid
								? "bg-neutrals-surface0 text-neutrals-text hover:text-vibrant-orchid cursor-pointer hover:opacity-90"
								: "bg-neutrals-surface0 text-neutrals-text opacity-40"
						}`}
					>
						Continue
					</button>
				</form>
			</div>
		</div>
	);
}

export default AuthForm;

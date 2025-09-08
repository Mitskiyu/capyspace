"use client";
import { google, logo } from "@/assets";
import { checkEmail, login, signUp } from "@/lib/auth";
import { zodResolver } from "@hookform/resolvers/zod";
import Image from "next/image";
import { useMemo, useState } from "react";
import { useForm } from "react-hook-form";
import { z } from "zod";

type State = "email" | "signup" | "login";

function createSchema(state: State) {
	const schema = z.object({
		email: z.email("Email is invalid"),
		password: z.string().min(1, "Password is required").optional(),
		confirm: z.string().min(1, "Please confirm your password").optional(),
	});

	if (state === "signup") {
		return schema
			.refine(
				(data) => {
					if (
						data.password &&
						data.password.length > 0 &&
						data.password.length < 8
					) {
						return false;
					}
					return true;
				},
				{
					message: "Password must be at least 8 characters",
					path: ["password"],
				},
			)
			.refine(
				(data) => {
					if (data.password && data.confirm) {
						return data.password === data.confirm;
					}
					return true;
				},
				{
					message: "Passwords do not match",
					path: ["confirm"],
				},
			);
	}

	return schema;
}

function AuthForm() {
	const [state, setState] = useState<State>("email");
	const schema = useMemo(() => createSchema(state), [state]);

	type Inputs = z.infer<typeof schema>;

	const {
		register,
		handleSubmit,
		resetField,
		setError,
		clearErrors,
		trigger,
		formState: { errors, isValid, isSubmitting },
	} = useForm<Inputs>({
		resolver: zodResolver(schema),
		mode: "onTouched",
	});

	async function onSubmit(data: Inputs) {
		if (state === "email") {
			const result = await checkEmail(data.email);

			if (!result.ok) {
				setError("email", { type: "manual", message: result.error });
				console.error(result.error);
				return;
			}

			setState(result.exists ? "login" : "signup");
			return;
		}

		if (state === "signup") {
			if (!data.password) return;

			const result = await signUp(data.email, data.password);

			if (!result.ok) {
				setError("password", { type: "manual", message: result.error });
				console.error(result.error);
				return;
			}

			// call login -> redirect!
		}

		if (state === "login") {
			if (!data.password) return;

			const result = await login(data.email, data.password);

			if (!result.ok) {
				setError("password", { type: "manual", message: result.error });
				console.error(result.error);
				return;
			}

			// redirect!
		}
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
					{(errors.email || errors.password || errors.confirm) && (
						<div className="text-vibrant-coral mt-2 text-sm">
							<span>
								{errors.email?.message ||
									errors.password?.message ||
									errors.confirm?.message}
							</span>
						</div>
					)}
					<input
						{...register("email")}
						onChange={(e) => {
							if (state !== "email") setState("email");
							clearErrors();
							resetField("password");
							resetField("confirm");
							register("email").onChange(e);
						}}
						type="text"
						placeholder="Enter your email"
						className={`bg-neutrals-surface0 focus:ring-vibrant-bloom h-10 w-11/12 rounded-2xl border-none px-4 py-2 text-base ring-1 transition-colors duration-300 ease-out placeholder:text-base focus:ring-2 focus:outline-none sm:h-12 sm:text-lg sm:placeholder:text-lg ${
							errors.email || errors.password || errors.confirm
								? "mt-2"
								: "mt-4"
						} ${state === "email" ? "ring-neutrals-overlay0" : ""} ${errors.email ? "ring-vibrant-coral" : ""}`}
					/>

					{state === "signup" && (
						<>
							<input
								{...register("password")}
								onChange={(e) => {
									register("confirm").onChange(e);
									setTimeout(() => {
										trigger(["password", "confirm"]);
									}, 0);
								}}
								type="password"
								placeholder="Enter your password"
								className={`bg-neutrals-surface0 focus:ring-vibrant-bloom ring-neutrals-overlay0 mt-2 h-10 w-11/12 rounded-2xl border-none px-4 py-2 text-base ring-1 transition-colors duration-300 ease-out placeholder:text-base focus:ring-2 focus:outline-none sm:h-12 sm:text-lg sm:placeholder:text-lg ${
									errors.password ? "ring-vibrant-coral ring-2" : ""
								}`}
							/>
							<input
								{...register("confirm")}
								onChange={(e) => {
									register("confirm").onChange(e);
									setTimeout(() => {
										trigger(["password", "confirm"]);
									}, 0);
								}}
								type="password"
								placeholder="Confirm your password"
								className={`bg-neutrals-surface0 focus:ring-vibrant-bloom ring-neutrals-overlay0 mt-1.5 h-10 w-11/12 rounded-2xl border-none px-4 py-2 text-base ring-1 transition-colors duration-300 ease-out placeholder:text-base focus:ring-2 focus:outline-none sm:h-12 sm:text-lg sm:placeholder:text-lg ${
									errors.confirm ? "ring-vibrant-coral ring-2" : ""
								}`}
							/>
						</>
					)}

					{state === "login" && (
						<input
							{...register("password")}
							type="password"
							placeholder="Enter your password"
							className={`bg-neutrals-surface0 focus:ring-vibrant-bloom ring-neutrals-overlay0 mt-2 h-10 w-11/12 rounded-2xl border-none px-4 py-2 text-base ring-1 transition-colors duration-300 ease-out placeholder:text-base focus:ring-2 focus:outline-none sm:h-12 sm:text-lg sm:placeholder:text-lg ${
								errors.password ? "ring-vibrant-coral ring-2" : ""
							}`}
						/>
					)}

					<button
						disabled={!isValid || isSubmitting}
						className={`mt-2 mb-4 flex h-10 w-11/12 items-center justify-center rounded-2xl text-base transition duration-200 sm:h-12 sm:text-lg ${
							!isValid
								? "bg-neutrals-surface0 text-neutrals-text cursor-default opacity-40"
								: isSubmitting
									? "bg-neutrals-surface0 text-neutrals-text cursor-default opacity-70"
									: "bg-neutrals-surface0 text-neutrals-text hover:text-vibrant-orchid cursor-pointer hover:opacity-90"
						}`}
					>
						{isSubmitting ? (
							<div
								className="bg-[conic-gradient(from_0deg,transparent,theme(colors.vibrant.lilac),theme(colors.vibrant.orchid),transparent)] size-4.5 animate-spin rounded-full"
								style={{
									mask: "radial-gradient(circle, transparent 6px, black 8px)",
									WebkitMask:
										"radial-gradient(circle, transparent 6px, black 8px)",
								}}
							></div>
						) : (
							<span>Continue</span>
						)}
					</button>
				</form>
			</div>
		</div>
	);
}

export default AuthForm;

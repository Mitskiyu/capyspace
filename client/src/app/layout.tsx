import type { Metadata } from "next";
import { Lexend } from "next/font/google";
import "./globals.css";

const lexend = Lexend({
	fallback: ["Helvetica"],
	subsets: ["latin"],
});

export const metadata: Metadata = {
	title: "Capyspace",
	description: "",
};

export default function RootLayout({
	children,
}: Readonly<{
	children: React.ReactNode;
}>) {
	return (
		<html lang="en">
			<body className={`${lexend.className} mx-auto px-4 antialiased`}>
				{children}
			</body>
		</html>
	);
}

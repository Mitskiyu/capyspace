import { ArrowsOutSimpleIcon, MinusIcon, XIcon } from "@phosphor-icons/react";
import { motion, useDragControls } from "motion/react";

function StickyNote() {
	const dragControls = useDragControls();

	return (
		<motion.div
			drag
			dragMomentum={false}
			dragElastic={0}
			dragListener={false}
			dragControls={dragControls}
			whileDrag={{
				scale: 1.05,
				boxShadow: "0px 10px 20px rgba(0,0,0,0.2)",
			}}
			className="max-w-80 rounded-xl sm:max-w-96"
		>
			<div className="bg-neutrals-base flex h-64 w-full flex-col rounded-xl p-2 sm:h-80 sm:p-4">
				<motion.div
					onPointerDown={(e) => dragControls.start(e)}
					className="group flex w-full cursor-move items-center justify-start space-x-1"
				>
					<div className="bg-vibrant-sunset flex size-4 cursor-pointer items-center justify-center rounded-full transition-colors">
						<XIcon
							size={12}
							weight="bold"
							className="text-black/70 opacity-0 transition-opacity duration-200 group-hover:opacity-100"
						/>
					</div>
					<div className="bg-vibrant-buttercup flex size-4 cursor-pointer items-center justify-center rounded-full transition-colors">
						<MinusIcon
							size={12}
							weight="bold"
							className="text-black/70 opacity-0 transition-opacity duration-200 group-hover:opacity-100"
						/>
					</div>
					<div className="bg-vibrant-mint flex size-4 cursor-pointer items-center justify-center rounded-full transition-colors">
						<ArrowsOutSimpleIcon
							size={12}
							weight="bold"
							className="rotate-90 text-black/70 opacity-0 transition-opacity duration-200 group-hover:opacity-100"
						/>
					</div>
				</motion.div>
				<hr className="from-vibrant-dream via-vibrant-cloud to-vibrant-dream mt-2 flex h-1 w-full rounded-xl border-0 bg-gradient-to-r" />
				<textarea
					className="text-neutrals-text placeholder-neutrals-subtext1/60 mt-2 flex-1 resize-none border-none bg-transparent text-base leading-relaxed outline-none sm:text-lg"
					placeholder="Start writing here..."
				/>
			</div>
		</motion.div>
	);
}

export default StickyNote;

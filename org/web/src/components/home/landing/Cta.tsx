import { Button } from 'components/shared/ui/button';
import { motion } from 'framer-motion';

export const CTA = () => {
    const variants = {
        offscreen: {
            y: 100,
            opacity: 0,
        },
        onscreen: {
            y: 0,
            opacity: 1,
            transition: {
                type: 'spring',
                mass: 0.5,
                damping: 30,
                duration: 0.8,
            },
        },
    };

    return (
        <motion.div
            initial="offscreen"
            whileInView="onscreen"
            viewport={{ once: true, amount: 0.3 }}
            variants={variants}
            className="flex flex-col items-center justify-center p-8 sm:p-12"
        >
            <section className="flex flex-col sm:flex-row items-center justify-between w-full max-w-4xl bg-gradient-to-r from-[#001020] to-[#00243E] text-white rounded-[20px] shadow-lg p-8 sm:p-12">
                {/* Left Content */}
                <div className="flex-1 flex flex-col items-start">
                    <p className="text-lg font-semibold tracking-wider bg-clip-text text-transparent bg-gradient-to-r from-blue-400 via-blue-300 to-blue-500">
                        It takes 1 minute
                    </p>
                    <h2 className="text-3xl md:text-4xl font-bold mt-4 leading-tight">
                        The faster, easier way to explore satellites
                    </h2>
                    <p className="text-base md:text-lg mt-4 max-w-lg text-gray-300">
                        Jump in today and discover how easy it is to track, explore, and learn about satellites orbiting Earth in real-time with Project 2112.
                    </p>
                </div>

                {/* Right Content */}
                <div className="flex items-center mt-6 sm:mt-0 sm:ml-10">
                    <Button size="xl" className="p-7 mt-6 text-xl z-10" variant="outlinePrimary">
                        Start exploring now
                    </Button>
                </div>
            </section>
        </motion.div>
    );
};

export default CTA;

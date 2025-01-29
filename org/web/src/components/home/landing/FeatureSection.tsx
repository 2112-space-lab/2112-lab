import { Button } from 'components/shared/ui/button';
import { motion } from 'framer-motion';

export const FeatureSection = ({ title, subtitle, keyPoints, buttonText, buttonColor }) => {
    const styles = {
        container: 'w-full flex flex-col items-center p-10 sm:p-16 rounded-2xl',
        title: 'text-[28px] md:text-[36px] font-bold text-center text-white',
        subtitle: 'mt-4 text-sm md:text-base text-gray-300 text-center max-w-3xl',
        keyPointContainer: 'grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-3 gap-8 mt-8 max-w-[1200px]',
        keyPointCard:
            'flex flex-col items-start p-6 bg-[#001A33] rounded-[15px] shadow-md hover:shadow-lg transition-shadow duration-300',
        keyPointTitle: 'text-xl font-semibold text-white',
        keyPointDescription: 'mt-2 text-sm text-gray-300',
        button:
            'mt-10 py-3 px-6 font-medium text-base rounded-lg shadow-md bg-gradient-to-r from-blue-500 to-blue-700 text-white hover:from-blue-600 hover:to-blue-800 transition-colors',
    };

    const variants = {
        offscreen: { opacity: 0, y: 100 },
        onscreen: { opacity: 1, y: 0, transition: { type: 'spring', mass: 0.5, damping: 30, duration: 0.8 } },
    };

    return (
        <motion.div
            initial="offscreen"
            whileInView="onscreen"
            viewport={{ once: true, amount: 0.3 }}
            variants={variants}
            className={styles.container}
        >
            <h1 className={styles.title}>{title}</h1>
            <p className={styles.subtitle}>{subtitle}</p>

            <div className={styles.keyPointContainer}>
                {keyPoints.map((point, index) => (
                    <motion.div
                        key={index}
                        className={styles.keyPointCard}
                        whileHover={{ scale: 1.05 }}
                    >
                        <h3 className={styles.keyPointTitle}>{point.title}</h3>
                        <p className={styles.keyPointDescription}>{point.description}</p>
                    </motion.div>
                ))}
            </div>

            <motion.button
                whileHover={{ scale: 1.05 }}
            >

                <Button size="xl" className="p-7 mt-6 text-xl z-10" variant="outlineSecondary">
                    {buttonText}
                </Button>
            </motion.button>
        </motion.div>
    );
};

export default FeatureSection;

import { motion } from 'framer-motion';
import { Rating } from './Rating'; // Adjust the path as necessary

export const Testimonials = () => {
    const testimonials = [
        {
            id: 1,
            content: "Project 2112 has completely transformed how we track and analyze satellite data.",
            name: "John Doe",
            title: "Aerospace Engineer at NASA",
            img: "https://picsum.photos/id/64/48/48",
            rating: 4.5,
        },
        {
            id: 2,
            content: "The best platform for real-time satellite tracking and insights.",
            name: "Jane Doe",
            title: "Astronomy Enthusiast",
            img: "https://picsum.photos/id/65/48/48",
            rating: 5,
        },
        {
            id: 3,
            content: "I've discovered hundreds of satellites and learned so much about orbital mechanics.",
            name: "Alice Doe",
            title: "Professor of Astrophysics",
            img: "https://picsum.photos/id/669/48/48",
            rating: 4,
        },
    ];

    const styles = {
        container: "sm:py-16 py-10 flex flex-col items-center text-white", // Added navy background for the container
        heading: "text-[28px] md:text-[36px] font-bold text-center text-white mb-6", // Adjusted font size and ensured text is white for contrast
        grid: "grid gap-6 sm:grid-cols-2 lg:grid-cols-3 mt-8 max-w-[1200px] px-4", // Maintained clean layout for grid with padding
        card: "flex flex-col p-6 rounded-[15px] shadow-md hover:shadow-lg transition-shadow duration-300 bg-gradient-to-r from-[#00243E] to-[#001A33]", // Navy-gradient background for cards
        text: "text-[16px] leading-[24px] text-gray-300 mb-6", // Text in gray for readability on dark background
        profile: "flex items-center mt-4",
        avatar: "w-12 h-12 rounded-full border-2 border-[#00476B]", // Navy border for avatars to match theme
        info: "ml-4",
        name: "font-semibold text-white text-[18px]", // Name in white for clear visibility
        title: "text-gray-400 text-[14px]", // Subtitle in lighter gray for subtlety
    };


    const variants = {
        offscreen: { opacity: 0, y: 100 },
        onscreen: {
            opacity: 1,
            y: 0,
            transition: { type: "spring", stiffness: 80, damping: 20, duration: 0.6 },
        },
    };

    return (
        <section className={styles.container}>
            <h2 className={styles.heading}>What People Are Saying About Us</h2>
            <motion.div
                initial="offscreen"
                whileInView="onscreen"
                viewport={{ once: true, amount: 0.3 }}
                variants={variants}
                className={styles.grid}
            >
                {testimonials.map(({ id, content, name, title, img, rating }) => (
                    <motion.div key={id} variants={variants} className={styles.card}>
                        {/* Star Rating */}
                        <Rating className="mb-4" rating={rating} maxRating={5} size="medium" />

                        {/* Testimonial Content */}
                        <p className={styles.text}>{content}</p>

                        <div className={styles.profile}>
                            <img src={img} alt={name} className={styles.avatar} />
                            <div className={styles.info}>
                                <h4 className={styles.name}>{name}</h4>
                                <p className={styles.title}>{title}</p>
                            </div>
                        </div>
                    </motion.div>
                ))}
            </motion.div>
        </section>
    );
};

export default Testimonials;

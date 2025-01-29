'use client';

import { useEffect, useState } from 'react';
import { useTheme } from 'next-themes';
import { MoonIcon, SunIcon } from 'lucide-react';
import { motion, AnimatePresence } from 'framer-motion';

const ThemeSwitch = () => {
  const [mounted, setMounted] = useState(false);
  const { theme, setTheme } = useTheme();

  const currentTheme =
    theme === 'system'
      ? window?.matchMedia('(prefers-color-scheme: dark)').matches
        ? 'dark'
        : 'light'
      : theme;

  const toggleTheme = () => {
    setTheme(currentTheme === 'dark' ? 'light' : 'dark');
  };

  useEffect(() => setMounted(true), []);

  if (!mounted) {
    return <div className="w-6 h-6"></div>;
  }

  return (
    <button
      aria-label="Toggle Dark Mode"
      onClick={toggleTheme}
      className="relative w-6 h-6 flex items-center justify-center"
    >
      <AnimatePresence exitBeforeEnter>
        {currentTheme === 'dark' ? (
          <motion.div
            key="dark"
            initial={{ opacity: 0, translateY: 10 }}
            animate={{ opacity: 1, translateY: 0 }}
            exit={{ opacity: 0, translateY: -10 }}
            transition={{ duration: 0.3, ease: 'easeInOut' }}
            className="absolute"
          >
            <MoonIcon className="w-6 h-6" />
          </motion.div>
        ) : (
          <motion.div
            key="light"
            initial={{ opacity: 0, translateY: 10 }}
            animate={{ opacity: 1, translateY: 0 }}
            exit={{ opacity: 0, translateY: -10 }}
            transition={{ duration: 0.3, ease: 'easeInOut' }}
            className="absolute"
          >
            <SunIcon className="w-6 h-6" />
          </motion.div>
        )}
      </AnimatePresence>
    </button>
  );
};

export default ThemeSwitch;

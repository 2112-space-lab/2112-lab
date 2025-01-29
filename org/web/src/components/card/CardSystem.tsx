import clsx from 'clsx';
import Image from 'components/shared/Image';
import { GlowBg } from 'components/shared/ui/glow-bg';

export const CardComponent = ({
    className,
    title,
    description,
    titleComponent,
    descriptionComponent,
    imageSrc,
    imageAlt = '',
    imagePosition = 'top',
    imageShadow = 'soft',
    withBackground = false,
    withBackgroundGlow = false,
    variant = 'primary',
    backgroundGlowVariant = 'primary',
}: {
    className?: string;
    title?: string | React.ReactNode;
    titleComponent?: React.ReactNode;
    description?: string | React.ReactNode;
    descriptionComponent?: React.ReactNode;
    imageSrc?: string;
    imageAlt?: string;
    imagePosition?: 'top' | 'left' | 'right';
    imageShadow?: 'none' | 'soft' | 'hard';
    withBackground?: boolean;
    withBackgroundGlow?: boolean;
    variant?: 'primary' | 'secondary';
    backgroundGlowVariant?: 'primary' | 'secondary';
}) => {
    return (
        <div
            className={clsx(
                'flex flex-col rounded-md overflow-hidden p-4 shadow-md relative',
                withBackground && variant === 'primary'
                    ? 'bg-primary-100/20 dark:bg-primary-900/10'
                    : '',
                withBackground && variant === 'secondary'
                    ? 'bg-secondary-100/20 dark:bg-secondary-900/10'
                    : '',
                className
            )}
        >
            {withBackgroundGlow && (
                <GlowBg
                    className="absolute w-full h-full z-0 opacity-50"
                    variant={backgroundGlowVariant}
                />
            )}

            {imageSrc && imagePosition === 'top' && (
                <Image
                    className={clsx(
                        'w-full rounded-md mb-4',
                        imageShadow === 'soft' && 'shadow-md',
                        imageShadow === 'hard' && 'hard-shadow'
                    )}
                    src={imageSrc}
                    alt={imageAlt}
                    width={320}
                    height={180}
                />
            )}

            <div className="z-10 flex flex-col">
                {title ? (
                    <h3 className="text-xl font-bold">{title}</h3>
                ) : (
                    titleComponent
                )}

                {description ? (
                    <p className="mt-2 text-sm text-gray-600 dark:text-gray-300">
                        {description}
                    </p>
                ) : (
                    descriptionComponent
                )}
            </div>
        </div>
    );
};

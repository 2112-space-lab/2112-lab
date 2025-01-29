import Footer from 'components/footer/FooterAuthDefault';
import NavLink from 'components/link/NavLink';
import { SignIn, Waitlist } from '@clerk/nextjs';
import Logo from '../landing/Logo';
import LandingHeader from '../landing/Header';

function LoginPage(props: {
    bgGradient?: string;
    bgImage?: string;
    logo?: JSX.Element;
    footerText?: string;
    learnMoreLink?: string;
    user?: { name: string };
}) {
    const {
        footerText = "Learn more about 2112 on",
        learnMoreLink = "https://2112.com/about",
        user,
    } = props;

    return (
        <div className="flex flex-col w-full min-h-screen">
            <header className="relative z-50">
                <LandingHeader
                    className="absolute top-0 w-full"
                    logo={<Logo className="h-9 w-auto" />}
                    logoDark={<Logo className="h-9 w-auto" />}
                />
            </header>
            <div className="relative flex bg-gradient-to-r from-[#001020] via-[#001530] to-[#000810]">
                <div className="mx-auto flex min-h-full w-full flex-col justify-start pt-12 md:max-w-[75%] lg:h-screen lg:max-w-[1013px] lg:px-8 lg:pt-0 xl:h-[100vh] xl:max-w-[1383px] xl:px-0 xl:pl-[70px]">
                    <div className="mb-auto flex flex-col pl-5 pr-5 md:pl-12 md:pr-0 lg:max-w-[48%] lg:pl-0 xl:max-w-full">
                        {!user ? (
                            <div className="flex items-center h-screen">
                                <SignIn routing="hash" withSignUp={false} fallbackRedirectUrl={"/admin/default"} />
                            </div>
                        ) : (
                            <NavLink href="/admin/default" className="flex items-center h-screen mt-0 w-max lg:pt-10">
                                <div className="mx-auto flex h-fit w-fit items-center cursor-pointer">
                                    <svg
                                        width="8"
                                        height="12"
                                        viewBox="0 0 8 12"
                                        fill="none"
                                        xmlns="http://www.w3.org/2000/svg"
                                    >
                                        <path
                                            d="M6.70994 2.11997L2.82994 5.99997L6.70994 9.87997C7.09994 10.27 7.09994 10.9 6.70994 11.29C6.31994 11.68 5.68994 11.68 5.29994 11.29L0.709941 6.69997C0.319941 6.30997 0.319941 5.67997 0.709941 5.28997L5.29994 0.699971C5.68994 0.309971 6.31994 0.309971 6.70994 0.699971C7.08994 1.08997 7.09994 1.72997 6.70994 2.11997V2.11997Z"
                                            fill="#A3AED0"
                                        />
                                    </svg>
                                    <p className="ml-3 text-sm text-gray-600">
                                        Continue as {user.name}
                                    </p>
                                </div>
                            </NavLink>

                        )}
                        <div className="absolute right-0 hidden h-full min-h-screen md:block lg:w-[49vw] 2xl:w-[44vw]">
                            <div
                                className={`absolute flex h-full w-full items-end justify-center bg-[#001A33] bg-cover bg-center lg:rounded-bl-[120px] xl:rounded-bl-[200px]`}
                            >
                                <div className="relative flex h-full w-full">
                                    <div className="absolute top-1/3 left-1/2 transform -translate-x-1/2 -translate-y-1/3 text-center">
                                        <div className="flex flex-col items-center justify-center h-full">
                                            <h2 className="text-3xl font-bold text-white mb-6">
                                                No account yet? <br></br>
                                            </h2>
                                            <Waitlist
                                                appearance={{
                                                    elements: {
                                                        card: 'shadow-lg rounded-lg p-6',
                                                        button: 'rounded-md px-4 py-2 hover:bg-blue-400 text-white font-semibold',
                                                    },
                                                }}
                                                afterJoinWaitlistUrl="/thank-you"
                                            />
                                        </div>
                                    </div>
                                </div>
                            </div>
                        </div>
                    </div>
                    <Footer />
                </div>
            </div>
        </div>
    );
}

export default LoginPage;

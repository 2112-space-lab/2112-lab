'use client';
import { useState } from 'react';
import Stepper from 'components/admin/contexts/Stepper';
import StepperControl from 'components/admin/contexts/StepperControl';
import { UseContextProvider } from 'components/admin/contexts/StepperContext';
import ContextDetails from 'components/admin/contexts/new-context/steps/ContextDetails';
import AdditionalSettings from 'components/admin/contexts/new-context/steps/AdditionalSettings';
import ReviewAndSubmit from 'components/admin/contexts/new-context/steps/ReviewAndSubmit';
import Card from 'components/card';

const NewContext = () => {
  const [currentStep, setCurrentStep] = useState(1);

  // Steps for context creation
  const steps = [
    { stepNo: 1, name: 'Context Details' },
    { stepNo: 2, name: 'Additional Settings' },
    { stepNo: 3, name: 'Review and Submit' },
  ];

  // Render the step content based on the current step
  const displayStep = (step: {
    stepNo: number;
    name: string;
    highlighted?: boolean;
    selected?: boolean;
    completed?: boolean;
  }) => {
    switch (step.stepNo) {
      case 1:
        return <ContextDetails />;
      case 2:
        return <AdditionalSettings />;
      case 3:
        return <ReviewAndSubmit />;
      default:
        return null;
    }
  };

  // Handle navigation between steps
  const handleClick = (direction: string) => {
    let newStep = currentStep;

    direction === 'next' ? newStep++ : newStep--;
    // Ensure the step is within bounds
    newStep > 0 && newStep <= steps.length && setCurrentStep(newStep);
  };

  return (
    <div className="mt-3 h-full w-full">
      <div className="h-[350px] w-full rounded-[20px] bg-gradient-to-br from-brand-400 to-brand-600 md:h-[390px]" />
      <div className="w-md:2/3 mx-auto h-full w-5/6 md:px-3  3xl:w-7/12">
        <div className="-mt-[280px] w-full pb-10 md:-mt-[240px] md:px-[70px]">
          <Stepper
            action={setCurrentStep}
            steps={steps}
            currentStep={currentStep}
          />
        </div>

        <Card extra={'h-full mx-auto pb-3'}>
          <div className="rounded-[20px]">
            <UseContextProvider>
              {displayStep(steps[currentStep - 1])}
            </UseContextProvider>
          </div>
          {/* Navigation buttons */}
          <StepperControl
            handleClick={handleClick}
            currentStep={currentStep}
            steps={steps}
          />
        </Card>
      </div>
    </div>
  );
};

export default NewContext;

import React from "react";

export const ReviewAndSubmit = () => {
    return (
        <div className="h-full w-full rounded-[20px] px-3 pt-7 md:px-8">
            {/* Header */}
            <h4 className="pt-[5px] text-xl font-bold text-navy-700 dark:text-white">
                Review and Submit
            </h4>

            {/* Content */}
            <div className="mt-4">
                <p className="text-sm text-gray-500">
                    Please review the details of your new context before submitting.
                </p>
                {/* Placeholder for showing context summary */}
                <div className="my-4 border rounded-lg p-4">
                    {/* Example of showing submitted data */}
                    <p><strong>Context Name:</strong> Sample Context</p>
                    <p><strong>Tenant ID:</strong> T12345</p>
                    <p><strong>Description:</strong> Sample description...</p>
                    <p><strong>Max Satellites:</strong> 10</p>
                    <p><strong>Max Tiles:</strong> 50</p>
                </div>

                {/* Submit Button */}
                <button className="mt-4 w-full rounded bg-brand-500 px-4 py-2 text-white">
                    Submit
                </button>
            </div>
        </div>
    );
};

export default ReviewAndSubmit;

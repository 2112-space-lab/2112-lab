import InputField from "components/fields/InputField";

export const ContextDetails = () => {
    return (
        <div className="h-full w-full rounded-[20px] px-3 pt-7 md:px-8">
            {/* Header */}
            <h4 className="pt-[5px] text-xl font-bold text-navy-700 dark:text-white">
                Context Details
            </h4>

            {/* Content */}
            <div className="mt-4">
                {/* Context Name */}
                <InputField
                    extra="mb-3"
                    label="Context Name"
                    placeholder="Enter context name"
                    id="contextName"
                    type="text"
                />

                {/* Tenant ID */}
                <InputField
                    extra="mb-3"
                    label="Tenant ID"
                    placeholder="Enter tenant ID"
                    id="tenantId"
                    type="text"
                />

                {/* Description */}
                <InputField
                    extra="mb-3"
                    label="Description"
                    placeholder="Enter a brief description (optional)"
                    id="description"
                    type="text"
                />
            </div>
        </div>
    );
};

export default ContextDetails;

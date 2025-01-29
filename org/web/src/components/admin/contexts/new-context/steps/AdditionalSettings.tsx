import InputField from "components/fields/InputField";

export const AdditionalSettings = () => {
    return (
        <div className="h-full w-full rounded-[20px] px-3 pt-7 md:px-8">
            {/* Header */}
            <h4 className="pt-[5px] text-xl font-bold text-navy-700 dark:text-white">
                Additional Settings
            </h4>

            {/* Content */}
            <div className="mt-4">
                {/* Max Satellite */}
                <InputField
                    extra="mb-3"
                    label="Max Satellites"
                    placeholder="Enter maximum number of satellites"
                    id="maxSatellite"
                    type="number"
                />

                {/* Max Tiles */}
                <InputField
                    extra="mb-3"
                    label="Max Tiles"
                    placeholder="Enter maximum number of tiles"
                    id="maxTiles"
                    type="number"
                />

                {/* Activated At */}
                <InputField
                    extra="mb-3"
                    label="Activated At"
                    placeholder="Select activation date (optional)"
                    id="activatedAt"
                    type="datetime-local"
                />
            </div>
        </div>
    );
};

export default AdditionalSettings;


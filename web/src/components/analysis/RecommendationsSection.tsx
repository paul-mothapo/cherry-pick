import { CheckCircle } from "lucide-react";

interface RecommendationsSectionProps {
  recommendations: string[];
}

export default function RecommendationsSection({
  recommendations,
}: RecommendationsSectionProps) {
  return (
    <div className="mt-6">
      <h4 className="font-medium text-eerie-black text-lg mb-3">
        Recommendations:
      </h4>
      <ul className="space-y-2">
        {recommendations.map((rec, idx) => (
          <li key={idx} className="flex items-start">
            <CheckCircle className="h-5 w-5 text-green-500 mr-2 mt-0.5 flex-shrink-0" />
            <span className="text-onyx">{rec}</span>
          </li>
        ))}
      </ul>
    </div>
  );
}

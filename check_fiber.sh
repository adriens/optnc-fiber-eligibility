#!/bin/bash
# check_fiber.sh - Vérifier l'éligibilité fibre pour plusieurs numéros

set -e

API_URL="${API_URL:-http://localhost:8080}"
OUTPUT_FILE="${OUTPUT_FILE:-fiber_check.csv}"

echo "🔍 Vérification d'éligibilité fibre OPT-NC"
echo "=========================================="
echo ""

# Create CSV header
echo "Numéro,Trouvé,Fibre Disponible,Statut,Contact" > "$OUTPUT_FILE"

# Array of phone numbers to check
phone_numbers=("257364" "286320")

# If arguments provided, use them instead
if [ $# -gt 0 ]; then
  phone_numbers=("$@")
fi

echo "📞 Numéros à vérifier: ${phone_numbers[*]}"
echo ""

for phone in "${phone_numbers[@]}"; do
  echo -n "Vérification $phone... "
  
  # Make API call
  response=$(http --body --ignore-stdin GET "$API_URL/api/v1/eligibility" phone=="$phone" 2>/dev/null)
  
  if [ $? -eq 0 ]; then
    # Parse JSON response
    numero=$(echo "$response" | jq -r '.data.phone_number // "N/A"')
    trouve=$(echo "$response" | jq -r '.data.found // false')
    disponible=$(echo "$response" | jq -r '.data.fiber.available // false')
    statut=$(echo "$response" | jq -r '.data.fiber.status // "unknown"')
    contact=$(echo "$response" | jq -r '.data.contact_phone // "N/A"')
    
    # Add to CSV
    echo "$numero,$trouve,$disponible,$statut,$contact" >> "$OUTPUT_FILE"
    
    # Display result
    if [ "$disponible" = "true" ]; then
      echo "✅ Fibre disponible"
    else
      echo "❌ Fibre non disponible ($statut)"
    fi
  else
    echo "⚠️  Erreur API"
    echo "$phone,error,error,error,N/A" >> "$OUTPUT_FILE"
  fi
done

echo ""
echo "📊 Résultats sauvegardés dans: $OUTPUT_FILE"
echo ""
cat "$OUTPUT_FILE" | column -t -s','

echo ""
echo "✅ Vérification terminée"

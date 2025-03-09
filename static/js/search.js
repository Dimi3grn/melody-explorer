// Fonctionnalité de recherche
document.addEventListener('DOMContentLoaded', function() {
    // Récupérer le formulaire de recherche
    const searchForm = document.querySelector('.search-container form');
    
    if (searchForm) {
        // Récupérer les cases à cocher de type
        const artistCheckbox = document.getElementById('type-artist');
        const albumCheckbox = document.getElementById('type-album');
        const trackCheckbox = document.getElementById('type-track');
        
        if (artistCheckbox && albumCheckbox && trackCheckbox) {
            // S'assurer qu'au moins une case est cochée à tout moment
            const typeCheckboxes = [artistCheckbox, albumCheckbox, trackCheckbox];
            
            typeCheckboxes.forEach(checkbox => {
                checkbox.addEventListener('change', function() {
                    // Si l'utilisateur essaie de tout décocher, garder la case actuelle cochée
                    if (!artistCheckbox.checked && !albumCheckbox.checked && !trackCheckbox.checked) {
                        this.checked = true;
                    }
                });
            });
            
            // Ajouter un écouteur d'événement pour la soumission du formulaire de recherche
            searchForm.addEventListener('submit', function(e) {
                const searchInput = this.querySelector('input[name="q"]');
                
                // Valider l'entrée de recherche
                if (searchInput.value.trim() === '') {
                    e.preventDefault();
                    alert('Veuillez saisir une requête de recherche');
                    return;
                }
                
                // S'assurer qu'au moins une case à cocher est cochée
                if (!artistCheckbox.checked && !albumCheckbox.checked && !trackCheckbox.checked) {
                    e.preventDefault();
                    alert('Veuillez sélectionner au moins un type de recherche (Artistes, Albums ou Pistes)');
                    // Par défaut, tous les types si aucun n'est sélectionné
                    artistCheckbox.checked = true;
                    albumCheckbox.checked = true;
                    trackCheckbox.checked = true;
                }
            });
        }
    }
});
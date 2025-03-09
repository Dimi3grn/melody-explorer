// Fonctionnalité de pagination
document.addEventListener('DOMContentLoaded', function() {
    // Récupérer les liens de pagination
    const paginationLinks = document.querySelectorAll('.pagination a');
    
    if (paginationLinks.length > 0) {
        // Ajouter un écouteur d'événement pour les liens de pagination
        paginationLinks.forEach(link => {
            link.addEventListener('click', function(e) {
                e.preventDefault();
                
                // Obtenir l'URL
                const url = new URL(this.href);
                const params = url.searchParams;
                
                // Mettre à jour le formulaire avec la nouvelle page
                const filterForm = document.getElementById('filter-form');
                
                if (filterForm) {
                    // Définir la valeur de l'entrée de page
                    let pageInput = filterForm.querySelector('input[name="page"]');
                    
                    if (!pageInput) {
                        pageInput = document.createElement('input');
                        pageInput.type = 'hidden';
                        pageInput.name = 'page';
                        filterForm.appendChild(pageInput);
                    }
                    
                    pageInput.value = params.get('page');
                    
                    // Soumettre le formulaire
                    filterForm.submit();
                } else {
                    // Pas de formulaire de filtre, simplement naviguer vers l'URL
                    window.location.href = this.href;
                }
            });
        });
    }
    
    // Gérer les changements de limite d'éléments
    const limitSelect = document.getElementById('limit');
    
    if (limitSelect) {
        limitSelect.addEventListener('change', function() {
            // Réinitialiser la pagination lorsque la limite change
            const pageInput = document.querySelector('input[name="page"]');
            if (pageInput) {
                pageInput.value = '1';
            }
            
            // Soumettre le formulaire
            const filterForm = document.getElementById('filter-form');
            if (filterForm) {
                filterForm.submit();
            }
        });
    }
});
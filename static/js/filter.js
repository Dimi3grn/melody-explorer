// Fonctionnalité de filtrage
document.addEventListener('DOMContentLoaded', function() {
    // Récupérer le formulaire de filtrage
    const filterForm = document.getElementById('filter-form');
    
    if (filterForm) {
        // Ajouter un écouteur d'événement pour les changements de filtre
        const filterInputs = filterForm.querySelectorAll('select');
        
        filterInputs.forEach(input => {
            input.addEventListener('change', function() {
                // Réinitialiser la pagination lorsque les filtres changent
                const pageInput = document.querySelector('input[name="page"]');
                if (pageInput) {
                    pageInput.value = '1';
                }
                
                // Soumettre le formulaire
                filterForm.submit();
            });
        });
        
        // Gérer le basculement des filtres sur mobile
        const filterToggle = document.getElementById('filter-toggle');
        const filterControls = document.getElementById('filter-controls');
        
        if (filterToggle && filterControls) {
            filterToggle.addEventListener('click', function() {
                filterControls.classList.toggle('active');
                
                if (filterControls.classList.contains('active')) {
                    filterToggle.innerText = 'Masquer les filtres';
                } else {
                    filterToggle.innerText = 'Afficher les filtres';
                }
            });
        }
    }
    
    // Gérer l'accordéon pour les questions fréquentes
    const faqItems = document.querySelectorAll('.faq-item h3');
    
    faqItems.forEach(item => {
        item.addEventListener('click', function() {
            const answer = this.nextElementSibling;
            
            // Basculer la visibilité de la réponse
            if (answer.style.display === 'none' || !answer.style.display) {
                answer.style.display = 'block';
                this.classList.add('active');
            } else {
                answer.style.display = 'none';
                this.classList.remove('active');
            }
        });
    });
    
    // Initialiser les éléments FAQ (masquer les réponses initialement)
    faqItems.forEach(item => {
        const answer = item.nextElementSibling;
        answer.style.display = 'none';
    });
});
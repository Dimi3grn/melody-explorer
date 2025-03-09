// Fonctionnalité des favoris
document.addEventListener('DOMContentLoaded', function() {
    // Récupérer tous les boutons de favoris
    const favoriteButtons = document.querySelectorAll('.btn-favorite');
    
    if (favoriteButtons.length > 0) {
        // Ajouter un écouteur d'événement pour les boutons de favoris
        favoriteButtons.forEach(button => {
            button.addEventListener('click', function() {
                // Obtenir les données de l'élément
                const id = this.getAttribute('data-id');
                const type = this.getAttribute('data-type');
                const name = this.getAttribute('data-name');
                const imageURL = this.getAttribute('data-image');
                
                // Vérifier si déjà en favori
                const isFavorite = this.classList.contains('active');
                
                if (isFavorite) {
                    // Supprimer des favoris
                    removeFavorite(id, type, this);
                } else {
                    // Ajouter aux favoris
                    addFavorite(id, type, name, imageURL, this);
                }
            });
        });
    }
    
    // Récupérer tous les boutons de suppression de favoris
    const removeFavoriteButtons = document.querySelectorAll('.btn-remove-favorite');
    
    if (removeFavoriteButtons.length > 0) {
        // Ajouter un écouteur d'événement pour les boutons de suppression de favoris
        removeFavoriteButtons.forEach(button => {
            button.addEventListener('click', function() {
                // Obtenir les données de l'élément
                const id = this.getAttribute('data-id');
                const type = this.getAttribute('data-type');
                
                // Supprimer des favoris
                removeFavoriteFromPage(id, type, this);
            });
        });
    }
    
    // Fonction pour ajouter un favori
    function addFavorite(id, type, name, imageURL, button) {
        if (!id || !type) {
            console.error('Données requises manquantes pour le favori:', { id, type, name });
            showNotification('Erreur lors de l\'ajout aux favoris: Données manquantes', 'error');
            return;
        }
        
        // Préparer les données de la requête
        const data = {
            id: id,
            type: type,
            name: name || 'Sans titre',
            image_url: imageURL || ''
        };
        
        console.log('Ajout aux favoris:', data);
        
        // Faire la requête API
        fetch('/api/favorites/add', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => {
            if (response.ok) {
                return response.json();
            }
            throw new Error('Échec de l\'ajout aux favoris');
        })
        .then(data => {
            if (data.success) {
                // Mettre à jour le bouton
                button.classList.add('active');
                const icon = button.querySelector('i');
                if (icon) {
                    icon.className = 'fas fa-heart';
                }
                
                // Mettre à jour le texte du bouton s'il en a
                if (button.innerText.includes('Add to Favorites')) {
                    button.innerText = '';
                    button.appendChild(icon);
                    button.appendChild(document.createTextNode('Remove from Favorites'));
                }
                
                // Afficher un message de succès
                showNotification('Ajouté aux favoris !', 'success');
            }
        })
        .catch(error => {
            console.error('Erreur:', error);
            showNotification('Échec de l\'ajout aux favoris. Veuillez réessayer.', 'error');
        });
    }
    
    // Fonction pour supprimer un favori
    function removeFavorite(id, type, button) {
        if (!id || !type) {
            console.error('Données requises manquantes pour la suppression du favori:', { id, type });
            showNotification('Erreur lors de la suppression des favoris: Données manquantes', 'error');
            return;
        }
        
        // Préparer les données de la requête
        const data = {
            id: id,
            type: type
        };
        
        console.log('Suppression du favori:', data);
        
        // Faire la requête API
        fetch('/api/favorites/remove', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        })
        .then(response => {
            if (response.ok) {
                return response.json();
            }
            throw new Error('Échec de la suppression des favoris');
        })
        .then(data => {
            if (data.success) {
                // Mettre à jour le bouton
                button.classList.remove('active');
                const icon = button.querySelector('i');
                if (icon) {
                    icon.className = 'far fa-heart';
                }
                
                // Mettre à jour le texte du bouton s'il en a
                if (button.innerText.includes('Remove from Favorites')) {
                    button.innerText = '';
                    button.appendChild(icon);
                    button.appendChild(document.createTextNode('Add to Favorites'));
                }
                
                // Afficher un message de succès
                showNotification('Supprimé des favoris !', 'success');
            }
        })
        .catch(error => {
            console.error('Erreur:', error);
            showNotification('Échec de la suppression des favoris. Veuillez réessayer.', 'error');
        });
    }
    
    // Fonction pour supprimer un favori de la page des favoris
    function removeFavoriteFromPage(id, type, button) {
        // D'abord supprimer du backend
        removeFavorite(id, type, button);
        
        // Puis supprimer de l'interface avec une animation
        const card = button.closest('.artist-card, .album-card, .track-card');
        if (card) {
            card.style.transition = 'opacity 0.3s ease';
            card.style.opacity = '0';
            
            setTimeout(() => {
                card.remove();
                
                // Vérifier s'il n'y a plus d'éléments de ce type
                const section = card.closest('.favorites-section');
                if (section) {
                    const remainingCards = section.querySelectorAll('.artist-card, .album-card, .track-card');
                    if (remainingCards.length === 0) {
                        section.remove();
                    }
                }
                
                // Vérifier s'il n'y a plus du tout de favoris
                const favoritesGrid = document.querySelectorAll('.favorites-grid');
                if (favoritesGrid.length === 0) {
                    const container = document.querySelector('.favorites-page .container');
                    if (container) {
                        const noFavorites = document.createElement('div');
                        noFavorites.className = 'no-favorites';
                        noFavorites.innerHTML = `
                            <p>Vous n'avez pas encore ajouté de favoris.</p>
                            <a href="/collection" class="btn btn-primary">Parcourir la collection</a>
                        `;
                        container.appendChild(noFavorites);
                    }
                }
            }, 300);
        }
    }
    
    // Fonction pour afficher une notification
    function showNotification(message, type) {
        // Vérifier si le conteneur de notification existe
        let notificationContainer = document.querySelector('.notification-container');
        
        if (!notificationContainer) {
            // Créer le conteneur de notification
            notificationContainer = document.createElement('div');
            notificationContainer.className = 'notification-container';
            document.body.appendChild(notificationContainer);
        }
        
        // Créer la notification
        const notification = document.createElement('div');
        notification.className = `notification ${type}`;
        notification.textContent = message;
        
        // Ajouter le bouton de fermeture
        const closeButton = document.createElement('button');
        closeButton.className = 'notification-close';
        closeButton.innerHTML = '&times;';
        closeButton.addEventListener('click', function() {
            notification.remove();
        });
        
        notification.appendChild(closeButton);
        
        // Ajouter la notification au conteneur
        notificationContainer.appendChild(notification);
        
        // Supprimer la notification après 3 secondes
        setTimeout(() => {
            notification.classList.add('fade-out');
            setTimeout(() => {
                notification.remove();
            }, 300);
        }, 3000);
    }
    
    // Ajouter les styles de notification s'ils ne sont pas déjà ajoutés
    if (!document.getElementById('notification-styles')) {
        const style = document.createElement('style');
        style.id = 'notification-styles';
        style.innerHTML = `
            .notification-container {
                position: fixed;
                top: 20px;
                right: 20px;
                z-index: 1000;
            }
            
            .notification {
                background-color: #fff;
                border-radius: 5px;
                box-shadow: 0 3px 10px rgba(0, 0, 0, 0.2);
                padding: 15px 20px;
                margin-bottom: 10px;
                position: relative;
                min-width: 200px;
                max-width: 300px;
                animation: slide-in 0.3s ease;
            }
            
            .notification.success {
                border-left: 5px solid #2ecc71;
            }
            
            .notification.error {
                border-left: 5px solid #e74c3c;
            }
            
            .notification.fade-out {
                animation: fade-out 0.3s ease forwards;
            }
            
            .notification-close {
                position: absolute;
                top: 5px;
                right: 5px;
                background: none;
                border: none;
                font-size: 1.2rem;
                cursor: pointer;
                color: #888;
            }
            
            @keyframes slide-in {
                from {
                    transform: translateX(100%);
                    opacity: 0;
                }
                to {
                    transform: translateX(0);
                    opacity: 1;
                }
            }
            
            @keyframes fade-out {
                from {
                    transform: translateX(0);
                    opacity: 1;
                }
                to {
                    transform: translateX(100%);
                    opacity: 0;
                }
            }
        `;
        
        document.head.appendChild(style);
    }
});
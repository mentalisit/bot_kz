package webapp

// HelperFunctions возвращает вспомогательные JavaScript функции
const HelperFunctions = `
    <!-- Вспомогательные функции -->
    <script>
        // Утилиты
        function escapeHtml(unsafe) {
            return unsafe
                .replace(/&/g, "&amp;")
                .replace(/</g, "&lt;")
                .replace(/>/g, "&gt;")
                .replace(/"/g, "&quot;")
                .replace(/'/g, "&#039;");
        }

        function showNotification(message, type = 'info') {
            if (window.Telegram && Telegram.WebApp) {
                Telegram.WebApp.showPopup({
                    title: type === 'error' ? 'Ошибка' : 'Успех',
                    message: message,
                    buttons: [{ type: 'ok' }]
                });
            } else {
                alert((type === 'error' ? '❌ ' : '✅ ') + message);
            }
        }

        function showError(message) {
            const errorDiv = document.getElementById('errorMessage');
            if (errorDiv) {
                errorDiv.textContent = message;
                errorDiv.style.display = 'block';
            }
        }

        function hideError() {
            const errorDiv = document.getElementById('errorMessage');
            if (errorDiv) {
                errorDiv.style.display = 'none';
            }
        }

        // Управление модальными окнами
        function showCreateRoleModal() {
            if (!currentUser) {
                showNotification('Не удалось определить пользователя. Пожалуйста, обновите страницу.', 'error');
                return;
            }
            const modal = document.getElementById('createRoleModal');
            if (modal) {
                modal.style.display = 'block';
            }
            hideError();
        }

        function hideCreateRoleModal() {
            const modal = document.getElementById('createRoleModal');
            if (modal) {
                modal.style.display = 'none';
            }
            const roleName = document.getElementById('roleName');
            const roleDescription = document.getElementById('roleDescription');
            if (roleName) roleName.value = '';
            if (roleDescription) roleDescription.value = '';
        }

        // Закрытие модального окна по клику вне его
        window.onclick = function(event) {
            const modal = document.getElementById('createRoleModal');
            if (event.target === modal) {
                hideCreateRoleModal();
            }
        }
    </script>`

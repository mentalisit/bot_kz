package webapp

// MainScript возвращает основную JavaScript логику
const MainScript = `
    <script>
        // Глобальные переменные
        let currentUser = null;
        let currentChatId = 0;
        let currentChatTitle = 'Общие роли';

        // Инициализация при загрузке страницы
        document.addEventListener('DOMContentLoaded', function() {
            console.log('DOM loaded, initializing WebApp...');
            initializeWebApp();
        });

        // Основная функция инициализации
        function initializeWebApp() {
            // Получаем параметры URL
            const urlParams = new URLSearchParams(window.location.search);
           // Получаем start параметр из tgWebAppStartParam
    const startParam = urlParams.get('tgWebAppStartParam');
    console.log('🎯 tgWebAppStartParam:', startParam);
    
    // Извлекаем chat ID из tgWebAppStartParam
    if (startParam && startParam.startsWith('chat')) {
        const chatIdStr = startParam.substring(4); // Отрезаем "chat"
        currentChatId = parseInt(chatIdStr) || 0;
        console.log('✅ Chat ID from tgWebAppStartParam:', currentChatId);
    }
    
    // Fallback: если не нашли в tgWebAppStartParam, пробуем chat_id
    if (currentChatId === 0) {
        currentChatId = parseInt(urlParams.get('chat_id')) || 0;
        console.log('✅ Chat ID from chat_id parameter:', currentChatId);
    }

            console.log('Initializing WebApp for chat ID:', currentChatId);

            // Инициализируем Telegram Web App
            if (window.Telegram && Telegram.WebApp) {
                const tg = Telegram.WebApp;
                tg.expand();
                tg.ready();
                
                console.log('Telegram WebApp initialized:', {
                    version: tg.version,
                    platform: tg.platform,
                    initData: tg.initData,
                    initDataUnsafe: tg.initDataUnsafe
                });

                // Инициализируем пользователя и чат
                initializeUser(tg);
                initializeChatInfo();
                
            } else {
                console.error('Telegram WebApp not available');
                showError('Telegram WebApp не доступен. Пожалуйста, откройте через Telegram.');
                initializeChatInfo();
                loadRoles(); // Все равно пытаемся загрузить роли
            }

            // Дополнительная проверка через 2 секунды
            setTimeout(function() {
                if (!currentUser) {
                    console.log('User still not initialized after 2 seconds');
                    initializeChatInfo();
                    loadRoles();
                }
            }, 2000);
        }

        // Инициализация информации о пользователе
        function initializeUser(tg) {
            if (tg.initDataUnsafe && tg.initDataUnsafe.user) {
                currentUser = {
                    id: tg.initDataUnsafe.user.id,
                    first_name: tg.initDataUnsafe.user.first_name || '',
                    last_name: tg.initDataUnsafe.user.last_name || '',
                    username: tg.initDataUnsafe.user.username || ''
                };
                
                console.log('User data received:', currentUser);
                
                updateUserInfo();
                sendUserDataToServer(currentUser, tg);
                loadRoles();
                
            } else {
                console.error('No user data found in initDataUnsafe:', tg.initDataUnsafe);
                showError('Не удалось получить данные пользователя');
                updateUserInfo();
                loadRoles();
            }
        }

        // Обновление информации о пользователе в интерфейсе
        function updateUserInfo() {
            const userInfoDiv = document.getElementById('userInfo');
            if (!userInfoDiv) return;
            
            if (currentUser) {
                userInfoDiv.innerHTML = 
                    '<p>Добро пожаловать, <strong>' + (currentUser.first_name || 'Пользователь') + '</strong>!</p>' +
                    '<p><small>ID: ' + currentUser.id + ' | @' + (currentUser.username || 'без username') + '</small></p>';
            } else {
                userInfoDiv.innerHTML = '<p>Не удалось получить данные пользователя</p>';
            }
        }

        // Инициализация информации о чате
        function initializeChatInfo() {
            const chatInfoDiv = document.getElementById('chatInfo');
            if (!chatInfoDiv) return;
            
            if (currentChatId === 0) {
                chatInfoDiv.innerHTML = '<strong>📱 Общие роли</strong><br><small>Роли доступные всем</small>';
                currentChatTitle = 'Общие роли';
            } else {
                chatInfoDiv.innerHTML = '<strong>👥 Роли группового чата</strong><br><small>ID: ' + currentChatId + '</small>';
                currentChatTitle = 'Групповой чат ID: ' + currentChatId;
            }
        }

        // Отправка данных пользователя на сервер
        async function sendUserDataToServer(userData, tg) {
            try {
                console.log('Sending user data to server:', userData);
                const response = await fetch('/api/user-info', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        user: userData,
                        initData: tg.initData,
                        chat_type: tg.initDataUnsafe.chat_type,
                        chat_id: currentChatId
                    })
                });
                
                if (response.ok) {
                    console.log('Данные пользователя отправлены на сервер');
                } else {
                    console.error('Ошибка отправки данных пользователя:', response.status);
                }
            } catch (error) {
                console.error('Ошибка отправки данных:', error);
            }
        }

        // Загрузка ролей
        async function loadRoles() {
            try {
                console.log('Loading roles for chat ID:', currentChatId);
                const userID = currentUser ? currentUser.id : 0;
                const response = await fetch('/api/roles?user_id=' + userID + '&chat_id=' + currentChatId);
                const result = await response.json();
                
                console.log('Roles API response:', result);
                
                if (result.status === 'success') {
                    renderRoles(result.data);
                } else {
                    showNotification('Ошибка загрузки ролей: ' + (result.message || ''), 'error');
                }
            } catch (error) {
                console.error('Ошибка загрузки ролей:', error);
                showNotification('Ошибка загрузки ролей', 'error');
            }
        }

        // Создание роли
        async function createRole() {
            const name = document.getElementById('roleName').value.trim();
            const description = document.getElementById('roleDescription').value.trim();

            if (!name) {
                showNotification('Введите название роли', 'error');
                return;
            }

            if (!currentUser) {
                showNotification('Не удалось определить пользователя. Пожалуйста, обновите страницу.', 'error');
                console.error('Current user is null when creating role');
                return;
            }

            console.log('Creating role:', { 
                name, 
                description, 
                created_by: currentUser.id,
                chat_id: currentChatId,
                chat_title: currentChatTitle
            });

            try {
                const response = await fetch('/api/create-role', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        name: name,
                        description: description,
                        created_by: currentUser.id,
                        chat_id: currentChatId,
                        chat_title: currentChatTitle
                    })
                });

                const result = await response.json();
                console.log('Create role response:', result);
                
                if (response.ok && result.status === 'success') {
                    hideCreateRoleModal();
                    loadRoles();
                    showNotification('Роль создана успешно!', 'success');
                } else {
                    showNotification(result.message || 'Ошибка создания роли', 'error');
                }
            } catch (error) {
                console.error('Ошибка создания роли:', error);
                showNotification('Ошибка создания роли: ' + error.message, 'error');
            }
        }

        // Подписка на роль
        async function subscribeToRole(roleId) {
            if (!currentUser) {
                showNotification('Не удалось определить пользователя', 'error');
                return;
            }

            try {
                const response = await fetch('/api/subscribe', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        user_id: currentUser.id,
                        role_id: roleId
                    })
                });

                if (response.ok) {
                    loadRoles();
                    showNotification('Подписка оформлена!', 'success');
                }
            } catch (error) {
                console.error('Ошибка подписки:', error);
                showNotification('Ошибка подписки', 'error');
            }
        }

        // Отписка от роли
        async function unsubscribeFromRole(roleId) {
            if (!currentUser) {
                showNotification('Не удалось определить пользователя', 'error');
                return;
            }

            try {
                const response = await fetch('/api/unsubscribe', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        user_id: currentUser.id,
                        role_id: roleId
                    })
                });

                if (response.ok) {
                    loadRoles();
                    showNotification('Отписка выполнена', 'success');
                }
            } catch (error) {
                console.error('Ошибка отписки:', error);
                showNotification('Ошибка отписки', 'error');
            }
        }

        // Удаление роли
        async function deleteRole(roleId) {
            if (!confirm('Вы уверены, что хотите удалить эту роль?')) {
                return;
            }

            if (!currentUser) {
                showNotification('Не удалось определить пользователя', 'error');
                return;
            }

            try {
                const response = await fetch('/api/delete-role', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({
                        role_id: roleId,
                        user_id: currentUser.id
                    })
                });

                if (response.ok) {
                    loadRoles();
                    showNotification('Роль удалена', 'success');
                } else {
                    showNotification('Не удалось удалить роль', 'error');
                }
            } catch (error) {
                console.error('Ошибка удаления роли:', error);
                showNotification('Ошибка удаления роли', 'error');
            }
        }

        // Рендеринг списка ролей
function renderRoles(roles) {
    const container = document.getElementById('rolesList');
    if (!container) return;
    
    if (roles.length === 0) {
        container.innerHTML = '<div style="text-align: center; color: var(--tg-theme-hint-color); margin: 40px 0; font-size: 16px;">Роли еще не созданы</div>';
        return;
    }

    let rolesHTML = '';
    for (let i = 0; i < roles.length; i++) {
        const role = roles[i];
        const isSubscribed = role.subscribed;
        const isOwner = currentUser && role.created_by === currentUser.id;
        
        rolesHTML += '<div class="role-card">' +
            '<div class="role-header">' +
            '<div class="role-name">' + escapeHtml(role.name) + '</div>' +
            '</div>' +
            '<div class="role-description">' + escapeHtml(role.description || 'Описание отсутствует') + '</div>' +
            '<div class="role-stats">' +
            '👥 Подписчиков: ' + role.subscribers.length +
            '</div>' +
            '<div class="action-buttons">';
        
        if (isSubscribed) {
            rolesHTML += '<button class="action-btn subscribed" onclick="unsubscribeFromRole(\'' + role.id + '\')">✅ Отписаться</button>';
        } else {
            rolesHTML += '<button class="action-btn not-subscribed" onclick="subscribeToRole(\'' + role.id + '\')">📝 Подписаться</button>';
        }
        
        if (isOwner) {
            rolesHTML += '<button class="action-btn delete-btn" onclick="deleteRole(\'' + role.id + '\')">🗑️ Удалить</button>';
        }
        
        rolesHTML += '</div></div>';
    }
    
    container.innerHTML = rolesHTML;
}
    </script>
</body>
</html>`

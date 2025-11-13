package webapp

// HTMLBody возвращает основную структуру HTML body с улучшенной разметкой
const HTMLBody = `<body>
    <div class="container">
        <div class="header">
            <h1>🎭 Управление ролями</h1>
            <div id="chatInfo" class="chat-info">
                <strong>Загрузка информации о чате...</strong>
                <small>Пожалуйста, подождите</small>
            </div>
            <div id="userInfo" class="user-info">Загрузка информации о пользователе...</div>
            <div id="errorMessage" class="error" style="display: none;"></div>
        </div>

        <button class="btn" onclick="showCreateRoleModal()">➕ Создать роль</button>
        <button class="btn" onclick="loadRoles()">🔄 Обновить список</button>

        <div id="rolesList"></div>

        <!-- Модальное окно создания роли -->
        <div id="createRoleModal" class="modal">
            <div class="modal-content">
                <h3>Создать новую роль</h3>
                <input type="text" id="roleName" class="form-input" placeholder="Название роли" required>
                <textarea id="roleDescription" class="form-input" placeholder="Описание роли (необязательно)" rows="3"></textarea>
                <div class="action-buttons">
                    <button class="btn" onclick="createRole()">Создать</button>
                    <button class="btn" style="background: var(--tg-theme-secondary-bg-color, #6c757d); color: var(--tg-theme-text-color, #ffffff);" onclick="hideCreateRoleModal()">Отмена</button>
                </div>
            </div>
        </div>
    </div>`

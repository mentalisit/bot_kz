package webapp

// HTMLHead возвращает HTML head с адаптивными стилями
const HTMLHead = `<!DOCTYPE html>
<html>
<head>
 <meta charset="UTF-8">
    <title>WebApp</title>
    <title>Управление ролями</title>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <script src="https://telegram.org/js/telegram-web-app.js"></script>
    <style>
        /* Базовые сбросы стилей */
        * {
            box-sizing: border-box;
            margin: 0;
            padding: 0;
        }
        
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Oxygen, Ubuntu, Cantarell, sans-serif;
            padding: 16px;
            margin: 0; 
            background: var(--tg-theme-bg-color, #ffffff);
            color: var(--tg-theme-text-color, #000000);
            line-height: 1.5;
            font-size: 16px;
            -webkit-font-smoothing: antialiased;
            -moz-osx-font-smoothing: grayscale;
        }
        
        .container { 
            max-width: 100%;
            margin: 0 auto;
            padding: 0 8px;
        }
        
        .header { 
            text-align: center; 
            margin-bottom: 24px;
            padding-bottom: 16px;
            border-bottom: 1px solid var(--tg-theme-section-separator-color, #e5e5e5);
        }
        
        .header h1 {
            font-size: clamp(20px, 6vw, 24px);
            font-weight: 600;
            margin-bottom: 12px;
            color: var(--tg-theme-text-color, #000000);
        }
        
        .btn { 
            display: block; 
            width: 100%; 
            padding: 14px 16px;
            margin: 12px 0; 
            background: var(--tg-theme-button-color, #2481cc); 
            color: var(--tg-theme-button-text-color, #ffffff); 
            border: none; 
            border-radius: 12px; 
            cursor: pointer; 
            text-align: center;
            font-size: 17px;
            font-weight: 500;
            transition: opacity 0.2s ease;
            -webkit-tap-highlight-color: transparent;
        }
        
        .btn:active {
            opacity: 0.8;
        }
        
        .role-card { 
            background: var(--tg-theme-secondary-bg-color, #f8f9fa); 
            padding: 20px;
            margin: 16px 0; 
            border-radius: 16px; 
            border: 1px solid var(--tg-theme-section-separator-color, #e9ecef);
        }
        
        .role-header {
            display: flex;
            justify-content: space-between;
            align-items: flex-start;
            margin-bottom: 12px;
        }
        
        .role-name {
            font-size: 18px;
            font-weight: 600;
            color: var(--tg-theme-text-color, #000000);
            line-height: 1.3;
            margin-right: 12px;
        }
        
        .role-description {
            font-size: 15px;
            color: var(--tg-theme-hint-color, #6c757d);
            margin-bottom: 16px;
            line-height: 1.4;
        }
        
        .role-stats {
            font-size: 14px;
            color: var(--tg-theme-hint-color, #6c757d);
            margin-bottom: 16px;
        }
        
        .modal {
            display: none;
            position: fixed;
            top: 0;
            left: 0;
            width: 100%;
            height: 100%;
            background: rgba(0,0,0,0.5);
            z-index: 1000;
            backdrop-filter: blur(4px);
        }
        
        .modal-content {
            position: absolute;
            top: 50%;
            left: 50%;
            transform: translate(-50%, -50%);
            background: var(--tg-theme-bg-color, #ffffff);
            padding: 24px;
            border-radius: 20px;
            width: 90%;
            max-width: 400px;
            max-height: 80vh;
            overflow-y: auto;
            box-shadow: 0 10px 30px rgba(0,0,0,0.2);
        }
        
        .modal-content h3 {
            font-size: 20px;
            font-weight: 600;
            margin-bottom: 20px;
            color: var(--tg-theme-text-color, #000000);
            text-align: center;
        }
        
        .form-input {
            width: 100%;
            padding: 16px;
            border: 1px solid var(--tg-theme-section-separator-color, #e0e0e0);
            border-radius: 12px;
            font-size: 16px;
            margin-bottom: 16px;
            background: var(--tg-theme-bg-color, #ffffff);
            color: var(--tg-theme-text-color, #000000);
            font-family: inherit;
            transition: border-color 0.2s ease;
        }
        
        .form-input:focus {
            outline: none;
            border-color: var(--tg-theme-button-color, #2481cc);
        }
        
        .form-input::placeholder {
            color: var(--tg-theme-hint-color, #6c757d);
        }
        
        .action-buttons {
            display: flex;
            gap: 12px;
            margin-top: 20px;
            flex-wrap: wrap;
        }
        
        .action-btn {
            flex: 1;
            min-width: 120px;
            padding: 12px 16px;
            border: none;
            border-radius: 10px;
            font-size: 15px;
            font-weight: 500;
            cursor: pointer;
            transition: all 0.2s ease;
            text-align: center;
            -webkit-tap-highlight-color: transparent;
        }
        
        .subscribed { 
            background: #28a745; 
            color: white; 
        }
        
        .not-subscribed { 
            background: var(--tg-theme-button-color, #2481cc); 
            color: white; 
        }
        
        .delete-btn { 
            background: #dc3545; 
            color: white; 
        }
        
        .error { 
            color: #dc3545; 
            text-align: center; 
            margin: 16px 0;
            font-size: 15px;
            padding: 12px;
            background: rgba(220, 53, 69, 0.1);
            border-radius: 8px;
        }
        
        .chat-info { 
            background: var(--tg-theme-secondary-bg-color, #f8f9fa); 
            padding: 16px; 
            border-radius: 12px; 
            margin-bottom: 20px;
            text-align: center;
            border: 1px solid var(--tg-theme-section-separator-color, #e9ecef);
        }
        
        .chat-info strong {
            font-size: 16px;
            font-weight: 600;
            color: var(--tg-theme-text-color, #000000);
            display: block;
            margin-bottom: 4px;
        }
        
        .chat-info small {
            font-size: 14px;
            color: var(--tg-theme-hint-color, #6c757d);
        }
        
        .user-info {
            margin-bottom: 20px;
            text-align: center;
        }
        
        .user-info p {
            margin-bottom: 8px;
            font-size: 15px;
            color: var(--tg-theme-text-color, #000000);
        }
        
        .user-info strong {
            font-weight: 600;
        }
        
        .user-info small {
            font-size: 13px;
            color: var(--tg-theme-hint-color, #6c757d);
        }
        
        /* Адаптивность для очень маленьких экранов */
        @media (max-width: 360px) {
            body {
                padding: 12px;
            }
            
            .container {
                padding: 0 4px;
            }
            
            .role-card {
                padding: 16px;
                margin: 12px 0;
            }
            
            .action-buttons {
                gap: 8px;
            }
            
            .action-btn {
                min-width: 100px;
                font-size: 14px;
                padding: 10px 12px;
            }
        }
        
        /* Поддержка темной темы */
        @media (prefers-color-scheme: dark) {
            body:not([style]) {
                background: #1a1a1a;
                color: #ffffff;
            }
        }
        
        /* Улучшение доступности */
        @media (prefers-reduced-motion: reduce) {
            * {
                transition: none !important;
            }
        }
        
        /* Улучшение для iOS */
        @supports (-webkit-touch-callout: none) {
            .btn, .action-btn {
                min-height: 44px;
            }
        }
    </style>
</head>`

-- 1. Создаем функцию, которая будет отправлять уведомление
CREATE OR REPLACE FUNCTION notify_config_change()
    RETURNS trigger AS $$
DECLARE
    channel_name text := 'config_updates';
    payload text;
BEGIN
    -- Формируем полезную нагрузку: "схема.таблица"
    payload := TG_TABLE_SCHEMA || '.' || TG_TABLE_NAME;

    -- Отправляем уведомление
    PERFORM pg_notify(channel_name, payload);

    -- В триггерах AFTER результат игнорируется, но принято возвращать NEW
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Удаляем старые триггеры, если они существовали (чтобы избежать дублей)
DROP TRIGGER IF EXISTS trg_bot_config_update ON rs_bot.config_rs;
DROP TRIGGER IF EXISTS trg_bridge_config_update ON rs_bot.bridge_config;
DROP TRIGGER IF EXISTS trg_kzbot_config_update ON kzbot.config;

-- 3. Вешаем триггер на таблицу rs_bot.config_rs
CREATE TRIGGER trg_bot_config_update
    AFTER INSERT OR UPDATE OR DELETE ON rs_bot.config_rs
    FOR EACH STATEMENT EXECUTE FUNCTION notify_config_change();

-- 4. Вешаем триггер на таблицу rs_bot.bridge_config
CREATE TRIGGER trg_bridge_config_update
    AFTER INSERT OR UPDATE OR DELETE ON rs_bot.bridge_config
    FOR EACH STATEMENT EXECUTE FUNCTION notify_config_change();

-- Вешаем триггер на таблицу kzbot.config
CREATE TRIGGER trg_kzbot_config_update
    AFTER INSERT OR UPDATE OR DELETE ON kzbot.config
    FOR EACH STATEMENT EXECUTE FUNCTION notify_config_change();

LISTEN config_updates;

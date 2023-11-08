export const log = (() => {
    const ignore = () => {};

    const _log = {
        ERROR:{color: '#EDAFA8', level: 1 },
        WARN: {color: '#E8FBD4', level: 2 },
        DEBUG:{color: '#4E8E71', level: 3 },
        INFO: {color: '#699BE0', level: 4 },
        TRACE:{color: '#E9E7FA', level: 5 },

        DEFAULT: { level: 4 },

        set level(level) {
            _log.error = level >= _log.ERROR.level ? (message) => console.error(`%c${message}`, `color: ${_log.ERROR.color};`) : ignore;
            _log.warn = level >= _log.WARN.level ? (message) => console.warn(`%c${message}`, `color: ${_log.WARN.color};`) : ignore;
            _log.debug = level >= _log.DEBUG.level ? (message) => console.debug(`%c${message}`, `color: ${_log.DEBUG.color};`) : ignore;
            _log.info = level >= _log.INFO.level ? (message) => console.info(`%c${message}`, `color: ${_log.INFO.color};`) : ignore;
            _log.trace = level >= _log.TRACE.level ? (message) => console.trace(`%c${message}`, `color: ${_log.TRACE.color};`) : ignore;
        },

        get level() {
            return _log._level;
        },
    };

    _log.level = _log.INFO.level;

    return _log;
})();

export default log;

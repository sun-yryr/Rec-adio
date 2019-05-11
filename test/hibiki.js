angular.module("env", []).constant("env", "production").constant("apiHost", "https://vcms-api.hibiki-radio.jp").constant("serviceStatusAPI", "https://s3-ap-northeast-1.amazonaws.com/hibiki-status/service_status.json").constant("meta", {
    robots: "index, follow",
    title: "響 - HiBiKi Radio Station -",
    description: "トレーディングカードゲームの株式会社ブシロードのスピンアウトカンパニー、インターネットラジオポータルサイト株式会社響　-HiBiKi-。カードゲームラジオは響でチェック！",
    keywords: ["株式会社響", "響", "HiBiKi", "~hibiki~", "ひびき", "hibiki-radio.jp", "ブシロード", "インターネットラジオ", "カードゲームラジオ", "ポータルサイト", "声優"]
}).constant("logicalog_settings", {
    host: "hibiki-track.logica.bz"
}).constant("ga", {
    id: "UA-68093026-1"
}),
function() {
    angular.module("hibikiWeb", ["ngAnimate", "ngCookies", "ngTouch", "ngSanitize", "ngMessages", "ngResource", "ui.router", "env", "angular-loading-bar", "ngDialog", "720kb.socialshare", "720kb.datepicker", "slick", "angular-cache", "infinite-scroll", "growlNotifications", "uuid", "angulartics", "angulartics.google.analytics"])
}
.call(this),
function() {
    angular.module("hibikiWeb").service("VoicePlayer", ["$q", "$rootScope", "$sce", "$compile", "$timeout", "$interval", "Api", "videoSettings", "Notification", "ProgramPlayer", "logicaLog", "User", "platform", "env", function($q, $rootScope, $sce, $compile, $timeout, $interval, Api, videoSettings, Notification, ProgramPlayer, logicaLog, User, platform, env) {
        var Voice;
        return Voice = function() {
            function Voice(type, voice, cb) {
                this.player_id = null,
                this.player = null,
                this.canplay = null,
                this.ended = null,
                this.type = type,
                this.voice = voice,
                this.loading = !1,
                this.firstplay = !1,
                this.init(type, voice, cb)
            }
            return Voice.prototype.init = function(type, voice, cb) {
                var deferred, self;
                return 0 === angular.element("#voice-player-" + type).length && angular.element("body").append(angular.element("<div id='voice-player-" + type + "'>")),
                self = this,
                deferred = $q.defer(),
                this.setResource(voice).then(function(url) {
                    var _settings, browserText, scope, userAgent, video_tag;
                    return self.player_id = self.getPlayerId(type),
                    scope = $rootScope.$new(),
                    scope.player_id = self.player_id,
                    scope.title = self.getName(type),
                    video_tag = $compile(videoSettings.template)(scope),
                    angular.element("#voice-player-" + type).append(video_tag),
                    platform().mobile ? _settings = videoSettings.options.mobile : (_settings = videoSettings.options.pc,
                    userAgent = window.navigator.userAgent.toLowerCase(),
                    browserText = "",
                    browserText = -1 !== userAgent.indexOf("msie") || -1 !== userAgent.indexOf("trident") ? "ie" : -1 !== userAgent.indexOf("edge") ? "edge" : -1 !== userAgent.indexOf("chrome") ? "chrome" : -1 !== userAgent.indexOf("safari") ? "safari" : -1 !== userAgent.indexOf("firefox") ? "firefox" : -1 !== userAgent.indexOf("opera") ? "opera" : "",
                    console.log(window.navigator.userAgent),
                    console.log(browserText),
                    "safari" !== browserText || -1 === userAgent.indexOf("version/8") && -1 === userAgent.indexOf("version/9") || (_settings.techOrder = ["flash"])),
                    $timeout(function() {
                        return self.player = new window.streaksplayer(self.player_id,_settings,function() {}
                        ),
                        self.setEventHandlers(),
                        angular.isFunction(cb) && cb(),
                        deferred.resolve()
                    }, 0)
                }),
                deferred.promise
            }
            ,
            Voice.prototype.setResource = function(voice) {
                return voice && (this.voice = voice),
                Api.video_url(this.voice.video.id).then(function(_this) {
                    return function(url) {
                        var deferred, video_tag;
                        return deferred = $q.defer(),
                        _this.url = url,
                        _this.player && _this.player.src() !== _this.url && (_this.player.src(_this.url),
                        video_tag = angular.element("#" + _this.player_id + " video"),
                        video_tag.length > 0 && video_tag.attr("title", _this.getName(_this.type)),
                        _this.setLog(!0)),
                        deferred.resolve(_this.url),
                        deferred.promise
                    }
                }(this))
            }
            ,
            Voice.prototype.setLog = function(update) {
                var log_params;
                return this.player ? (log_params = {
                    category: "play",
                    action: "play",
                    video_id: this.voice.video.id,
                    vuid: logicaLog.uuid(),
                    dsid: logicaLog.session_id(),
                    env: env,
                    token: "U7JkK5Cd6NWi"
                },
                "daily" === this.type ? (log_params.label = "daily_voice",
                log_params.daily_voice_id = this.voice.id) : (log_params.label = "instantvoice",
                log_params.instant_voice_id = this.voice.id),
                User.loggedIn() && (log_params.user_id = User.get().id),
                update ? this.player.tracker().updateOptions({
                    logParams: log_params
                }) : this.player.tracker({
                    host: logicaLog.host(),
                    logParams: log_params,
                    span: 30
                })) : void 0
            }
            ,
            Voice.prototype.setEventHandlers = function() {
                var self;
                return self = this,
                this.player.off("load"),
                this.player.off("loaded"),
                this.player.off("loadedmetadata"),
                this.player.off("timeupdate"),
                this.player.off("play"),
                this.player.off("pause"),
                this.player.off("volumechange"),
                this.player.off("resize"),
                this.player.off("error"),
                this.player.off("firstplay"),
                this.player.off("canplay"),
                this.player.off("durationchange"),
                this.player.off("ended"),
                this.player.off("fullscreenchange"),
                this.player.on("firstplay", function() {
                    return self.firstplay = !0
                }),
                this.player.on("playing", function() {
                    return angular.isFunction(self.canplay) ? (self.canplay(),
                    self.loading = !1) : void 0
                }),
                this.player.on("canplay", function() {
                    return angular.isFunction(self.canplay) ? (self.canplay(),
                    self.loading = !1) : void 0
                }),
                this.player.on("pause", function() {
                    return angular.isFunction(self.ended) ? self.ended() : void 0
                }),
                this.player.on("ended", function() {
                    return angular.isFunction(self.ended) ? self.ended() : void 0
                }),
                this.player.on("error", function(e) {
                    return console.log(e),
                    Notification.alert("エラーが発生しました。"),
                    console.log("TODO: voice player error handling")
                }),
                this.setLog()
            }
            ,
            Voice.prototype.getPlayerId = function(type) {
                return "voice-lgcp-player-" + type
            }
            ,
            Voice.prototype.getName = function(type) {
                return "daily" === type ? this.voice.name : this.voice.actor.name + " - " + this.voice.description
            }
            ,
            Voice.prototype.playWithFade = function(canplay, ended) {
                var fadeIn, fadeOut, playCall, self;
                return self = this,
                this.canplay = canplay,
                this.ended = ended,
                fadeOut = function() {
                    var deferred, replay;
                    return deferred = $q.defer(),
                    replay = !1,
                    ProgramPlayer.player && (ProgramPlayer.paused() || ProgramPlayer.pause(),
                    ProgramPlayer.ctrl = !1),
                    deferred.resolve(replay),
                    deferred.promise
                }
                ,
                playCall = function(replay) {
                    var deferred;
                    return deferred = $q.defer(),
                    $timeout(function() {
                        return self.player.on("ended", function() {
                            return deferred.resolve(replay)
                        }),
                        self.play()
                    }, 100),
                    deferred.promise
                }
                ,
                fadeIn = function(replay) {
                    var deferred;
                    return deferred = $q.defer(),
                    ProgramPlayer.player && (replay && ProgramPlayer.paused() ? ProgramPlayer.fadein() : ProgramPlayer.ctrl = !0),
                    deferred.resolve(),
                    deferred.promise
                }
                ,
                $q.when().then(fadeOut).then(playCall).then(fadeIn)
            }
            ,
            Voice.prototype.load = function() {
                return this.player ? this.player.load() : void 0
            }
            ,
            Voice.prototype.play = function() {
                var playInterval;
                return playInterval = null,
                this.player ? (this.player.play(),
                playInterval = $interval(function(_this) {
                    return function() {
                        return _this.player && _this.player.paused() ? _this.player.play() : $interval.cancel(playInterval)
                    }
                }(this), 500, 10)) : void 0
            }
            ,
            Voice.prototype.pause = function() {
                return this.player ? (angular.isFunction(this.ended) && $timeout(function(_this) {
                    return function() {
                        return _this.ended()
                    }
                }(this), 0),
                this.player.pause()) : void 0
            }
            ,
            Voice.prototype.stop = function() {
                return this.player ? (angular.isFunction(this.ended) && this.ended(),
                this.player.dispose(),
                $timeout(function(_this) {
                    return function() {
                        return _this.player_id = null,
                        _this.player = null
                    }
                }(this), 100)) : void 0
            }
            ,
            Voice.prototype.paused = function() {
                return this.player ? this.player.paused() : !0
            }
            ,
            Voice
        }(),
        {
            voice: {},
            init: function(type, voice) {
                return this.voice[type] ? void 0 : this.voice[type] = new Voice(type,voice)
            },
            setResource: function(type, voice) {
                return this.voice[type] ? this.voice[type].setResource(voice) : void 0
            },
            playWithFade: function(type, canplay, ended) {
                return this.voice[type].playWithFade(canplay, ended)
            },
            play: function(type) {
                return this.voice[type].play()
            },
            pause: function(type) {
                return this.voice[type].paused() ? void 0 : this.voice[type].pause()
            },
            paused: function(type) {
                return this.voice[type].paused()
            },
            pauseAll: function() {
                return _.forEach(this.voice, function(v) {
                    return v.paused() ? void 0 : v.pause()
                })
            },
            load: function(type) {
                return this.voice[type] ? this.voice[type].load() : void 0
            },
            loading: function(type) {
                return this.voice[type].loading
            },
            firstplay: function(type) {
                return this.voice[type] ? this.voice[type].firstplay : void 0
            }
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").service("User", ["$q", "$cookies", "Api", "moment", "prefectures", "ngDialog", function($q, $cookies, Api, moment, prefectures, ngDialog) {
        var __user, _user_data, daysInMonth, user_data_key;
        return user_data_key = "Hibiki-user-token",
        _user_data = $cookies.getObject(user_data_key),
        angular.isUndefined(_user_data) && (_user_data = {
            id: null,
            email: null,
            name: null,
            birth_date: null,
            prefecture_id: null,
            gender: null,
            allow_mail_flg: null,
            access_token: null,
            expired_at: null
        }),
        daysInMonth = function(year, month) {
            var days, last_day;
            return year || (year = moment().subtract(20, "year").year()),
            month || (month = moment().month() + 1),
            last_day = moment(new Date(year,month - 1)).daysInMonth(),
            days = [],
            _.times(last_day, function(n) {
                return days.push(n + 1)
            }),
            days
        }
        ,
        __user = {
            model: {
                email: "",
                name: "",
                password: "",
                password_confirm: "",
                birth_date: "",
                birth_date_parts: {
                    year: moment().subtract(20, "year").year(),
                    month: moment().month() + 1,
                    day: moment().date()
                },
                prefecture_id: "",
                allow_mail_flg: "true",
                gender: "1"
            },
            initModel: function() {
                return this.model.email = "",
                this.model.name = "",
                this.model.password = "",
                this.model.password_confirm = "",
                this.model.birth_date = "",
                this.model.birth_date_parts = {
                    year: moment().subtract(20, "year").year(),
                    month: moment().month() + 1,
                    day: moment().date()
                },
                this.model.prefecture_id = "",
                this.model.allow_mail_flg = "true",
                this.model.gender = "1"
            },
            formParts: {
                birth_date: {
                    year: function() {
                        var now_year, old_year, years;
                        return now_year = moment().year(),
                        old_year = moment().subtract(100, "year").year(),
                        years = [],
                        _.times(now_year - old_year, function(n) {
                            return years.push(old_year + n)
                        }),
                        years
                    }(),
                    month: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12],
                    day: function() {
                        return daysInMonth()
                    }()
                },
                prefectures: prefectures
            },
            updateDaysInMonth: function() {
                return this.formParts.birth_date.day = daysInMonth(this.model.birth_date_parts.year, this.model.birth_date_parts.month)
            },
            set: function(data) {
                return angular.isDefined(data.id) && (_user_data.id = data.id),
                angular.isDefined(data.email) && (_user_data.email = data.email),
                angular.isDefined(data.name) && (_user_data.name = data.name),
                angular.isDefined(data.birth_date) && (_user_data.birth_date = data.birth_date),
                angular.isDefined(data.prefecture_id) && (_user_data.prefecture_id = data.prefecture_id),
                angular.isDefined(data.gender) && (_user_data.gender = data.gender),
                angular.isDefined(data.allow_mail_flg) && (_user_data.allow_mail_flg = data.allow_mail_flg),
                angular.isDefined(data.access_token) && (_user_data.access_token = data.access_token),
                angular.isDefined(data.expired_at) && (_user_data.expired_at = data.expired_at),
                $cookies.putObject(user_data_key, {
                    id: _user_data.id,
                    name: _user_data.name,
                    access_token: _user_data.access_token,
                    expired_at: _user_data.expired_at
                }, {
                    expires: _user_data.expired_at
                })
            },
            get: function() {
                return moment(new Date(_user_data.expired_at)).isBefore(new Date) && this.destroy(),
                _user_data
            },
            loggedIn: function() {
                return !!_user_data.id
            },
            limited: function(event) {
                return this.loggedIn() ? !1 : (event && (event.preventDefault(),
                event.stopPropagation()),
                ngDialog.open({
                    template: "app/components/user/limited.html",
                    trapFocus: !1,
                    preserveFocus: !1
                }),
                !0)
            },
            profile: function() {
                var deferred, self;
                return self = this,
                deferred = $q.defer(),
                Api.user_profile().then(function(result) {
                    return self.set(result.data),
                    deferred.resolve(result.data)
                }, function(e) {
                    return deferred.reject(e)
                }),
                deferred.promise
            },
            login: function(data) {
                var deferred, self;
                return deferred = $q.defer(),
                self = this,
                Api.login(data).then(function(result) {
                    return self.set(result.data),
                    Api.user_profile()
                }).then(function(result) {
                    return self.set(result.data),
                    deferred.resolve(result.data)
                })["catch"](function(e) {
                    return deferred.reject(e)
                }),
                deferred.promise
            },
            logout: function() {
                var deferred, self;
                return deferred = $q.defer(),
                self = this,
                Api.logout().then(function(result) {
                    return self.destroy(),
                    deferred.resolve()
                }, function(e) {
                    return deferred.reject(e)
                }),
                deferred.promise
            },
            destroy: function() {
                return $cookies.remove(user_data_key)
            },
            register: function() {
                var deferred, post_data, self;
                return deferred = $q.defer(),
                post_data = {
                    email: this.model.email,
                    name: this.model.name,
                    password: this.model.password,
                    prefecture_id: this.model.prefecture_id,
                    birth_date: this.model.birth_date,
                    allow_mail_flg: this.model.allow_mail_flg,
                    gender: this.model.gender
                },
                self = this,
                Api.register(post_data).then(function(result) {
                    return deferred.resolve(result.data)
                })["catch"](function(e) {
                    return deferred.reject(e)
                }),
                deferred.promise
            },
            update: function() {
                var deferred, post_data, self;
                return deferred = $q.defer(),
                post_data = {
                    email: this.model.email,
                    name: this.model.name,
                    prefecture_id: this.model.prefecture_id,
                    birth_date: this.model.birth_date,
                    allow_mail_flg: this.model.allow_mail_flg,
                    gender: this.model.gender
                },
                self = this,
                Api.user_update(post_data).then(function(result) {
                    return self.set(result.data),
                    Api.user_profile()
                }).then(function(result) {
                    return self.set(result.data),
                    deferred.resolve(result.data)
                })["catch"](function(e) {
                    return deferred.reject(e)
                }),
                deferred.promise
            },
            password_update: function(password, new_password) {
                var deferred, post_data, self;
                return deferred = $q.defer(),
                post_data = {
                    password: password,
                    new_password: new_password
                },
                self = this,
                Api.password_update(post_data).then(function(result) {
                    return deferred.resolve(result.data)
                })["catch"](function(e) {
                    return deferred.reject(e)
                }),
                deferred.promise
            },
            password_reminder: function(email) {
                var deferred, self;
                return deferred = $q.defer(),
                self = this,
                Api.password_reminder(email).then(function(result) {
                    return deferred.resolve(result.data)
                })["catch"](function(e) {
                    return deferred.reject(e)
                }),
                deferred.promise
            },
            withdraw: function() {
                var deferred, self;
                return deferred = $q.defer(),
                self = this,
                Api.withdraw().then(function(result) {
                    return self.destroy(),
                    deferred.resolve()
                })["catch"](function(e) {
                    return deferred.reject(e)
                }),
                deferred.promise
            }
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").directive("tooltip", ["$window", "$compile", "$interval", "$sce", function($window, $compile, $interval, $sce) {
        var default_options;
        return default_options = {
            showTrigger: "mouseover",
            hideTrigger: "mouseleave",
            position: "top",
            "class": ""
        },
        {
            restrict: "A",
            scope: {},
            link: function($scope, $element, $attributes) {
                var _template, _tooltip, content, element, hideTooltip, hideTriggers, setPosition, showTooltip, showTriggers, title, tooltip_class, view_position;
                return title = $attributes.tooltipTitle || $attributes.title || "",
                content = $attributes.tooltipContent || "",
                showTriggers = $attributes.tooltipShowTrigger || default_options.showTrigger,
                hideTriggers = $attributes.tooltipHideTrigger || default_options.hideTrigger,
                view_position = $attributes.tooltipPosition || default_options.position,
                tooltip_class = $attributes.tooltipClass || default_options["class"],
                element = angular.element($element[0]),
                $scope.title = title,
                $scope.content = content,
                $scope["class"] = tooltip_class,
                _template = '<div class="tooltip">\n  <div class="tooltip-title">{{title}}</div>\n  <div class="tooltip-content">{{content}}</div>\n  <span class="caret"></span>\n</div>',
                _tooltip = $compile(_template)($scope),
                _tooltip.addClass(view_position),
                tooltip_class.length > 0 && _tooltip.addClass(tooltip_class),
                setPosition = function() {
                    var element_height, element_width, left, offleft, offsetLeft, offsetTop, tooltip_height, tooltip_width, top;
                    return _tooltip.removeClass("top"),
                    _tooltip.removeClass("left"),
                    _tooltip.removeClass("right"),
                    _tooltip.removeClass("bottom"),
                    offsetTop = element.get(0).getBoundingClientRect().top + $window.scrollY,
                    isNaN(offsetTop) && (offsetTop = element.get(0).getBoundingClientRect().top + $window.pageYOffset),
                    offsetLeft = element.get(0).getBoundingClientRect().left + $window.scrollX,
                    isNaN(offleft) && (offleft = element.get(0).getBoundingClientRect().left + $window.pageXOffset),
                    element_width = element.get(0).offsetWidth,
                    element_height = element.get(0).ofssetHeight,
                    tooltip_width = _tooltip.get(0).offsetWidth,
                    tooltip_height = _tooltip.get(0).offsetHeight,
                    "right" === view_position ? (top = offsetTop + element_height / 2 - tooltip_height / 2,
                    left = offsetLeft + element_width + tooltip_width,
                    _tooltip.addClass("right")) : "left" === view_position ? (top = offsetTop + height / 2 - tooltip_height / 2,
                    left = offsetLeft - (tooltip_width + 10),
                    _tooltip.addClass("left")) : "top" === view_position ? (top = offsetTop - 10 - tooltip_height,
                    left = offsetLeft + element_width / 2 - tooltip_width / 2,
                    _tooltip.addClass("top")) : "bottom" === view_position && (top = offsetTop + element_height + 10,
                    left = offsetLeft + element_width / 2 - tooltip_width / 2,
                    _tooltip.addClass("bottom")),
                    _tooltip.css({
                        top: top,
                        left: left
                    })
                }
                ,
                showTooltip = function() {
                    return setPosition(),
                    _tooltip.addClass("active")
                }
                ,
                hideTooltip = function() {
                    return _tooltip.removeClass("active"),
                    _tooltip.css({
                        top: "",
                        left: ""
                    })
                }
                ,
                element.on("mouseover", function() {
                    return showTooltip()
                }),
                element.on("mouseleave", function() {
                    return hideTooltip()
                }),
                $attributes.tooltipTitle && $attributes.$observe("tooltipTitle", function(val) {
                    return $scope.title = val
                }),
                $attributes.title && $attributes.$observe("title", function(val) {
                    return $scope.title = val
                }),
                $attributes.content && $attributes.$observe("tooltipContent", function(val) {
                    return $scope.content = val
                }),
                $scope.$on("$destroy", function() {
                    return _tooltip.remove()
                }),
                angular.element("body").append(_tooltip)
            }
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").service("ProgramQueue", ["platform", "ProgramPlayer", "Notification", function(platform, ProgramPlayer, Notification) {
        var __queue;
        return __queue = {
            queue: [],
            isEmpty: function() {
                return 0 === this.queue.length
            },
            set: function(list) {
                return this.queue = list
            },
            push: function(program) {
                return this.queue.push(program)
            },
            pop: function() {
                var program;
                return this.isEmpty() ? null : (program = this.queue.shift(),
                program.episode && program.episode.video ? program : this.pop())
            },
            clear: function() {
                return this.queue = []
            },
            play: function() {
                var _q, program;
                return this.isEmpty() ? (Notification.alert("再生可能な番組がありません"),
                !1) : (program = this.pop(),
                program ? (ProgramPlayer.player || ProgramPlayer.init(program),
                _q = this,
                ProgramPlayer.setResource(program).then(function() {
                    return _q.isEmpty() || (ProgramPlayer.ended_callback = function() {
                        return _q.play()
                    }
                    ),
                    ProgramPlayer.auto_play = !0,
                    ProgramPlayer.play()
                })) : void 0)
            }
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").service("ProgramPlayer", ["$rootScope", "$sce", "$compile", "$timeout", "$interval", "$q", "Api", "videoSettings", "Notification", "logicaLog", "User", "platform", "env", function($rootScope, $sce, $compile, $timeout, $interval, $q, Api, videoSettings, Notification, logicaLog, User, platform, env) {
        var __player;
        return __player = {
            player_id: "program-player",
            ctrl: !1,
            loading: !0,
            standby: !1,
            firstplay: !1,
            live: !1,
            duration: 0,
            time: 0,
            pos: 0,
            volume_handle: !1,
            seek_handle: !1,
            player: null,
            url: "",
            backstage: !1,
            program_name: "",
            program_image: null,
            program_id: null,
            episode_id: null,
            episode_name: "",
            chapter_view: !1,
            chapter: [],
            chapter_pos: 0,
            current_chapter: null,
            screen_view: !1,
            media_type: "",
            fullscreen: !1,
            ended_callback: null,
            auto_play: !1,
            video: null,
            reset: function() {
                return this.ctrl = !1,
                this.loading = !0,
                this.standby = !1,
                this.firstplay = !1,
                this.live = !1,
                this.duration = 0,
                this.time = 0,
                this.pos = 0,
                this.volume_handle = !1,
                this.seek_handle = !1,
                this.screen_view = !1,
                this.chapter_view = !1,
                this.media_type = "",
                this.ended_callback = null,
                this.auto_play = !1,
                this.apply()
            },
            init: function(program, backstage) {
                var deferred, scope, self, video_tag;
                return this.player_ctrl = angular.element("#player-ctrl"),
                this.program_screen_container = angular.element("#program-screen-container"),
                this.seekBar = angular.element("#player-ctrl .seek-bar"),
                this.fullscreen_btn = angular.element("#player-ctrl .fullscreen-btn"),
                this.volume_icon = angular.element(".volume .volume-icon i"),
                this.volume_level = angular.element(".volume .volume-slider-level"),
                this.container = angular.element("#program-player-container"),
                0 === this.container.length && (angular.element("#program-screen").append(angular.element('<div id="program-player-container">')),
                this.container = angular.element("#program-player-container")),
                scope = $rootScope.$new(),
                video_tag = $compile(videoSettings.programPlayerTemplate)(scope),
                this.container.append(video_tag),
                self = this,
                deferred = $q.defer(),
                this.setResource(program, backstage).then(function() {}),
                deferred.promise
            },
            setResource: function(program, backstage) {
                return backstage ? (this.video = program.episode.additional_video,
                this.backstage = !0) : (this.video = program.episode.video,
                this.backstage = !1),
                Api.video_url(this.video.id).then(function(_this) {
                    return function(url) {
                        var _settings, browserText, deferred, first_chapter, userAgent, video_tag, video_tag_container, video_title;
                        return deferred = $q.defer(),
                        _this.url = url,
                        _this.program_id = program.id,
                        _this.episode_id = program.episode.id,
                        _this.video_id = _this.video.id,
                        _this.program_name = program.name,
                        _this.program_image = {
                            pc: program.pc_image_url,
                            sp: program.sp_image_url
                        },
                        _this.backstage ? (_this.episode_name = program.episode.name + "の楽屋裏",
                        _this.chapter = []) : (_this.episode_name = program.episode.name,
                        _this.chapter = _.sortBy(program.episode.chapters, function(x) {
                            return x.start_time
                        })),
                        _this.current_chapter = null,
                        _this.fullscreen_btn.removeClass("video").removeClass("audio"),
                        _this.container.removeClass("video").removeClass("audio"),
                        _this.player || (platform().mobile ? _settings = videoSettings.options.mobile : (_settings = videoSettings.options.pc,
                        userAgent = window.navigator.userAgent.toLowerCase(),
                        browserText = "",
                        browserText = -1 !== userAgent.indexOf("msie") || -1 !== userAgent.indexOf("trident") ? "ie" : -1 !== userAgent.indexOf("edge") ? "edge" : -1 !== userAgent.indexOf("chrome") ? "chrome" : -1 !== userAgent.indexOf("safari") ? "safari" : -1 !== userAgent.indexOf("firefox") ? "firefox" : -1 !== userAgent.indexOf("opera") ? "opera" : "",
                        "safari" !== browserText || -1 === userAgent.indexOf("version/8") && -1 === userAgent.indexOf("version/9") || (_settings.techOrder = ["flash"])),
                        _this.player = new window.streaksplayer("program-player",_settings),
                        _this.player.controls(!1),
                        _this.setEventHandlers(),
                        platform().mobile || _this.volumeUpdate(.5)),
                        _this.player && _this.player.src() !== _this.url && (_this.reset(),
                        _this.player.loadVideo({
                            sources: [{
                                src: _this.url,
                                type: "application/x-mpegURL"
                            }]
                        }),
                        video_tag_container = angular.element("#" + _this.player_id),
                        video_tag_container.length > 0 && (video_title = _this.program_name + " - " + _this.episode_name,
                        video_tag_container.attr("title", video_title),
                        video_tag = video_tag_container.children().find("video"),
                        video_tag.length > 0 && video_tag.attr("title", video_title)),
                        _this.setLog(!0)),
                        1 === _this.video.media_type ? (_this.media_type = "audio",
                        _this.fullscreen_btn.addClass("audio"),
                        _this.container.addClass("audio")) : 2 === _this.video.media_type && (_this.media_type = "video",
                        _this.fullscreen_btn.addClass("video"),
                        _this.container.addClass("video")),
                        first_chapter = _.first(_this.chapter),
                        first_chapter && 0 !== first_chapter.start_time && _this.chapter.unshift({
                            start_time: 0
                        }),
                        _this.setScreenImage(),
                        deferred.resolve(),
                        deferred.promise
                    }
                }(this))
            },
            setScreenImage: function() {
                var pc_image_url, sp_image_url, url;
                return pc_image_url = this.current_chapter && this.current_chapter.pc_image_url && this.current_chapter.pc_image_url.length > 0 ? this.current_chapter.pc_image_url : this.program_image.pc,
                sp_image_url = this.current_chapter && this.current_chapter.sp_image_url && this.current_chapter.sp_image_url.length > 0 ? this.current_chapter.sp_image_url : this.program_image.sp,
                url = this.fullscreen ? pc_image_url : sp_image_url,
                angular.element("#program-screen").css("background-image", "url(" + url + ")")
            },
            enable: function() {
                return this.url && this.standby
            },
            play: function() {
                var i, maxLoopCount, playInterval;
                return playInterval = null,
                this.player ? (this.player.play(),
                maxLoopCount = 10,
                i = 0,
                playInterval = $interval(function(_this) {
                    return function() {
                        return i >= maxLoopCount ? void 0 : (i++,
                        _this.player && _this.player.paused() ? _this.player.play() : i = maxLoopCount)
                    }
                }(this), 500, maxLoopCount)) : void 0
            },
            pause: function() {
                this.player && this.player.pause()
            },
            fadeout: function() {
                var fadeoutInterval, initial_volume;
                if (this.player)
                    return initial_volume = this.player.volume(),
                    fadeoutInterval = $interval(function(_this) {
                        return function() {
                            var _volume;
                            return _volume = _this.player.volume(),
                            0 === _volume ? ($interval.cancel(fadeoutInterval),
                            _this.player.pause(),
                            _this.ctrl = !1,
                            _this.volumeUpdate(initial_volume)) : _this.volumeUpdate(_volume - .1)
                        }
                    }(this), 200)
            },
            fadein: function() {
                var initial_volume;
                if (this.player)
                    return initial_volume = this.player.volume(),
                    this.ctrl = !0
            },
            playToggle: function() {
                this.player.paused() ? this.play() : this.pause()
            },
            paused: function() {
                if (!this.player)
                    return !0;
                try {
                    return this.player.paused()
                } catch (error) {}
            },
            rewind: function(sec) {
                var now, time;
                now = this.player.currentTime(),
                time = now - sec,
                0 > time && (time = 0),
                this.jumpTo(time)
            },
            jumpTo: function(time, force) {
                null == force && (force = !1),
                this.force && this.player.seeking() || this.player.currentTime(time)
            },
            jumpToChapter: function(chapter, adjust) {
                var time;
                if (!this.player.seeking())
                    return logicaLog.send({
                        category: "tap",
                        action: "jump_chapter",
                        label: this.program_id,
                        value: chapter.start_time
                    }),
                    time = adjust ? chapter.start_time + adjust : chapter.start_time,
                    this.jumpTo(time)
            },
            hasNext: function() {
                return this.chapter.length > this.chapter_pos + 1
            },
            hasPrev: function() {
                return this.time > 0
            },
            next: function() {
                var chapter;
                this.hasNext() && (chapter = this.chapter[this.chapter_pos + 1],
                this.jumpToChapter(chapter, .1))
            },
            prev: function() {
                var chapter;
                this.hasPrev() && (chapter = this.chapter[this.chapter_pos],
                null != chapter ? this.jumpToChapter(chapter, -.1) : this.jumpTo(0))
            },
            chapterPos: function(c) {
                return {
                    left: c.start_time / this.duration * 100 + "%"
                }
            },
            isChapterActive: function(time) {
                return this.time >= time
            },
            onChapterChange: function(chapter) {
                var target;
                return !this.standby || this.paused() || this.screen_view && this.chapter_view || (target = angular.element(".seek-bar .chapter-point").not(".active").first(),
                target.trigger("mouseover"),
                $timeout(function() {
                    return target.trigger("mouseleave")
                }, 3e3)),
                this.current_chapter = chapter,
                this.setScreenImage(),
                this.apply()
            },
            seekBarUpdate: function() {
                var _time, before_chapter_pos, chapter, is_live, tech_duration, time_format;
                this.player && (this.duration = this.player.duration(),
                _time = this.player.currentTime(),
                tech_duration = 0,
                is_live = function(_this) {
                    return function() {
                        var live, timerange;
                        return live = !1,
                        Infinity === _this.duration && (live = !0),
                        "Html5" === _this.player.techName_ && (timerange = _this.player.seekable(),
                        0 === timerange.length && (live = !0)),
                        live
                    }
                }(this),
                is_live() ? (this.live = !0,
                this.player_ctrl.addClass("live"),
                this.program_screen_container.addClass("live"),
                this.time = 0 > _time ? 0 : _time,
                this.video.dvr_flg ? (this.seekBar.find(".duration").text("Live"),
                this.seekBar.find(".time").text(moment.unix(this.time).format(this.timeFormat(this.time))),
                this.seekBar.find(".seek").css("width", "100%")) : this.seekBar.hide()) : (this.live = !1,
                this.player_ctrl.removeClass("live"),
                this.program_screen_container.removeClass("live"),
                this.time = 0 > _time ? 0 : _time > this.duration ? this.duration : _time,
                this.pos = this.time / this.duration * 100,
                this.seekBar.show(),
                chapter = _.last(_.filter(this.chapter, function(_this) {
                    return function(c) {
                        return _this.time - c.start_time >= 0
                    }
                }(this))),
                _.isUndefined(chapter) ? this.chapter_pos = 0 : (before_chapter_pos = this.chapter_pos,
                this.chapter_pos = _.indexOf(_.pluck(this.chapter, "start_time"), chapter.start_time),
                (null === this.current_chapter && 0 === chapter.start_time || before_chapter_pos !== this.chapter_pos) && this.onChapterChange(chapter)),
                time_format = this.timeFormat(this.duration),
                this.seekBar.find(".duration").text(moment.unix(this.duration).utc().format(time_format)),
                this.seekBar.find(".time").text(moment.unix(this.time).utc().format(time_format)),
                this.seekBar.find(".seek").css("width", this.pos + "%")))
            },
            seekHandle: function(e) {
                var _pos;
                return this.seek_handle = !0,
                this.seek_bar_rect = angular.element(".seek-bar .seek-container").get(0).getBoundingClientRect(),
                _pos = (e.pageX - this.seek_bar_rect.left) / this.seek_bar_rect.width * 10 / 10,
                this.jumpTo(this.duration * _pos),
                this.player.pause()
            },
            volumeUpdate: function(_volume) {
                _volume > 1 ? _volume = 1 : 0 > _volume && (_volume = 0),
                this.player.volume(_volume),
                this.volume_icon.removeClass("fa-volume-off").removeClass("fa-volume-down").removeClass("fa-volume-up"),
                0 === _volume ? this.volume_icon.addClass("fa-volume-off") : .5 >= _volume ? this.volume_icon.addClass("fa-volume-down") : 1 >= _volume && this.volume_icon.addClass("fa-volume-up"),
                this.volume_level.css("width", 100 * _volume + "%")
            },
            volumeHandle: function(e) {
                var _volume;
                return this.volume_handle = !0,
                this.volume_slider_rect = angular.element(".volume .volume-slider").get(0).getBoundingClientRect(),
                _volume = Math.round((e.pageX - this.volume_slider_rect.left) / this.volume_slider_rect.width * 10) / 10,
                this.volumeUpdate(_volume)
            },
            screen: function(on_or_off) {
                return this.screen_view = on_or_off
            },
            screenToggle: function() {
                return this.screen_view = !this.screen_view
            },
            fullScreen: function(on_or_off) {
                var d, screen_off, screen_on, self;
                return d = $q.defer(),
                self = this,
                screen_off = function() {
                    var deferred;
                    return deferred = $q.defer(),
                    self.screen_view = !1,
                    self.ctrl = !1,
                    $timeout(function() {
                        return deferred.resolve()
                    }, 300),
                    deferred.promise
                }
                ,
                screen_on = function() {
                    return $timeout(function() {
                        return self.screen_view = !0,
                        self.ctrl = !0
                    }, 200)
                }
                ,
                on_or_off ? screen_off().then(function() {
                    var deferred;
                    return deferred = $q.defer(),
                    self.fullscreen = on_or_off,
                    angular.element("#player-ctrl").addClass("fullscreen"),
                    angular.element("#program-screen-container").addClass("fullscreen"),
                    self.setScreenImage(),
                    $timeout(function() {
                        return $(window).trigger("resize")
                    }, 500),
                    deferred.resolve(),
                    deferred.promise
                }).then(screen_on).then(function() {
                    return d.resolve()
                }) : (self.fullscreen = on_or_off,
                screen_off().then(function() {
                    var deferred;
                    return deferred = $q.defer(),
                    angular.element("#player-ctrl").removeClass("fullscreen"),
                    angular.element("#program-screen-container").removeClass("fullscreen"),
                    self.setScreenImage(),
                    $timeout(function() {
                        return $(window).trigger("resize")
                    }, 500),
                    deferred.resolve(),
                    deferred.promise
                }).then(screen_on).then(function() {
                    return d.resolve()
                })),
                d.promise
            },
            fullscreenVideo: function() {
                return this.screen_view = !0,
                this.player.requestFullscreen()
            },
            clearHandle: function() {
                return this.volume_handle = !1,
                this.volume_slider_rect = null,
                this.seek_handle = !1,
                this.seek_bar_rect = null
            },
            apply: function() {
                return $rootScope.$$phase || $rootScope.$apply()
            },
            clear: function() {
                var deferred;
                return deferred = $q.defer(),
                this.fullScreen(!1).then(function(_this) {
                    return function() {
                        return _this.standby = !1,
                        _this.ctrl = !1,
                        _this.apply(),
                        deferred.resolve()
                    }
                }(this)),
                deferred.promise
            },
            timeFormat: function(time) {
                return time > 3600 ? "HH:mm:ss" : "mm:ss"
            },
            setLog: function(update) {
                var log_params;
                return this.player ? (log_params = {
                    category: "play",
                    action: "play",
                    label: this.backstage ? "episode_appendix" : "episode",
                    program_id: this.program_id,
                    episode_id: this.episode_id,
                    video_id: this.video_id,
                    vuid: logicaLog.uuid(),
                    dsid: logicaLog.session_id(),
                    env: env,
                    token: "U7JkK5Cd6NWi"
                },
                User.loggedIn() && (log_params.user_id = User.get().id),
                update ? this.player.tracker().updateOptions({
                    logParams: log_params
                }) : this.player.tracker({
                    host: logicaLog.host(),
                    logParams: log_params,
                    span: 30
                })) : void 0
            },
            setEventHandlers: function() {
                var self;
                return self = this,
                angular.element(document).on("mousemove", function(event) {
                    var _pos, _volume;
                    return self.volume_handle && self.volume_slider_rect && (_volume = Math.round((event.pageX - self.volume_slider_rect.left) / self.volume_slider_rect.width * 10) / 10,
                    self.volumeUpdate(_volume)),
                    self.seek_handle && self.seek_bar_rect && (_pos = (event.pageX - self.seek_bar_rect.left) / self.seek_bar_rect.width * 10 / 10,
                    _pos >= 1 ? _pos = .99 : 0 > _pos && (_pos = 0),
                    self.jumpTo(self.duration * _pos, !0),
                    !self.player.paused()) ? self.player.pause() : void 0
                }).on("mouseup", function() {
                    return self.seek_handle && self.play(),
                    self.clearHandle()
                }),
                this.player.on("firstplay", function() {
                    return self.screen_view = !0,
                    self.firstplay = !0
                }),
                this.player.on("canplay", function() {}),
                this.player.on("play", function() {
                    return self.ctrl = !0,
                    self.loading = !1,
                    self.chapter_pos = 0
                }),
                this.player.on("durationchange", function() {
                    return self.seekBarUpdate(),
                    self.volumeUpdate(self.player.volume()),
                    self.standby = !0
                }),
                this.player.on("timeupdate", function() {
                    return self.seekBarUpdate()
                }),
                this.player.on("ended", function() {
                    return self.apply(),
                    self.player.pause(),
                    self.player.controls(!1),
                    angular.isFunction(self.ended_callback) && self.ended_callback(),
                    !self.video.live_flg || self.video.replay_flg && self.video.dvr_flg ? void 0 : $timeout(function() {
                        return self.clear().then(function() {
                            return self.player.dispose(),
                            self.player = null
                        })
                    }, 200)
                }),
                this.player.on("fullscreenchange", function() {
                    return self.is_fullscreen_video ? $timeout(function() {
                        return self.is_fullscreen_video = !1,
                        self.player.controls(!1)
                    }, 100) : $timeout(function() {
                        return self.is_fullscreen_video = !0,
                        self.player.controls(!0)
                    }, 100)
                }),
                this.player.on("error", function(e) {
                    return console.log(e),
                    self.clear(),
                    Notification.alert("エラーが発生しました。")
                }),
                this.setLog()
            }
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").directive("programPlayerCtrl", function() {
        var ProgramPlayerController, directive;
        return ProgramPlayerController = function($rootScope, $scope, $window, $timeout, moment, platform, ProgramPlayer) {
            var focuscheck_type, live_watch, pgp;
            pgp = this,
            pgp.player = ProgramPlayer,
            pgp.platform = platform(),
            platform().mobile && (focuscheck_type = ["text", "email", "password"],
            angular.element(document).on("DOMFocusIn", function(e) {
                return e.target && ~focuscheck_type.indexOf(e.target.type) && ProgramPlayer.ctrl ? ProgramPlayer.ctrl = !1 : void 0
            }).on("DOMFocusOut", function(e) {
                return e.target && ~focuscheck_type.indexOf(e.target.type) && !ProgramPlayer.ctrl && ProgramPlayer.firstplay ? ProgramPlayer.ctrl = !0 : void 0
            })),
            live_watch = function() {
                var mb;
                if (ProgramPlayer.ctrl) {
                    if (pgp.platform.mobile) {
                        if ("ios" === pgp.platform.os)
                            return;
                        mb = ProgramPlayer.live ? "50px" : "70px"
                    } else
                        mb = "50px";
                    return angular.element("#layout .view-container, #menu").css("margin-bottom", mb)
                }
            }
            ,
            $scope.$watch(function() {
                return ProgramPlayer.ctrl
            }, live_watch),
            $scope.$watch(function() {
                return ProgramPlayer.live
            }, live_watch)
        }
        ,
        directive = {
            restrict: "E",
            templateUrl: "app/components/program-player/program-player.html",
            replace: !1,
            scope: {},
            controller: ProgramPlayerController,
            controllerAs: "pgp"
        }
    })
}
.call(this),
function() {
    angular.module("hibikiWeb").directive("programList", function() {
        var ProgramListController, directive;
        return ProgramListController = function($rootScope, $scope, $timeout, Api, moment, $cookies, $window, platform) {
            var copyrights, fixProgramBoxHeight, is_mobile, pl, updateCopyright;
            is_mobile = platform().mobile,
            pl = this,
            pl.$scope = $scope,
            pl.loaded = !1,
            pl.busy = !1,
            pl.page = 1,
            pl.limit = 8,
            pl.now_week = moment().locale("en").format("ddd").toLowerCase(),
            ("sat" === pl.now_week || "sun" === pl.now_week) && (pl.now_week = "satsun"),
            pl.isMonday = function() {
                return "mon" === pl.now_week
            }
            ,
            pl.isTuesday = function() {
                return "tue" === pl.now_week
            }
            ,
            pl.isWednesday = function() {
                return "wed" === pl.now_week
            }
            ,
            pl.isThursday = function() {
                return "thu" === pl.now_week
            }
            ,
            pl.isFriday = function() {
                return "fri" === pl.now_week
            }
            ,
            pl.isSaturdayOrSunday = function() {
                return "satsun" === pl.now_week || "sat" === pl.now_week || "sun" === pl.now_week
            }
            ,
            pl.isAllday = function() {
                return "all" === pl.now_week
            }
            ,
            copyrights = "",
            is_mobile ? (pl.name_length = 10,
            pl.latest_episode_name_length = 9,
            pl.cast_length = 12) : (pl.name_length = 32,
            pl.latest_episode_name_length = 15,
            pl.cast_length = 35),
            pl.dayOfWeekStr = function(day_of_week) {
                var str;
                switch (str = "",
                day_of_week) {
                case "mon":
                    str = is_mobile ? "月" : "月曜日";
                    break;
                case "tue":
                    str = is_mobile ? "火" : "火曜日";
                    break;
                case "wed":
                    str = is_mobile ? "水" : "水曜日";
                    break;
                case "thu":
                    str = is_mobile ? "木" : "木曜日";
                    break;
                case "fri":
                    str = is_mobile ? "金" : "金曜日";
                    break;
                case "satsun":
                    str = is_mobile ? "土・日" : "土・日曜日";
                    break;
                case "all":
                    str = "全て"
                }
                return str
            }
            ,
            pl.showPrograms = function(day_of_week) {
                return delete window.google_persistent_state_async,
                pl.busy = !0,
                pl.page = 1,
                angular.isUndefined(day_of_week) && (day_of_week = $cookies.selected_day_of_week ? $cookies.selected_day_of_week : pl.now_week),
                pl.now_week = $cookies.selected_day_of_week = $rootScope.selected_day_of_week = day_of_week,
                Api.programs({
                    day_of_week: day_of_week,
                    limit: pl.limit,
                    page: pl.page
                }).then(function(programs) {
                    return pl.loaded = !0,
                    pl.selected_programs = [],
                    $timeout(function() {
                        return pl.selected_programs = programs.data,
                        pl.selected_programs.length > 1 && pl.selected_programs.splice(1, 0, 1),
                        pl.selected_programs.length > 4 && pl.selected_programs.splice(4, 0, 4),
                        pl.selected_programs.length > 7 && pl.selected_programs.splice(7, 0, 7),
                        updateCopyright(),
                        fixProgramBoxHeight(),
                        pl.busy = programs.data.length < pl.limit ? !0 : !1,
                        !$scope.$$phase || $scope.$apply()
                    }, 200)
                })
            }
            ,
            pl.loadMore = function() {
                return pl.busy ? void 0 : (pl.busy = !0,
                pl.page += 1,
                Api.programs({
                    day_of_week: pl.now_week,
                    limit: pl.limit,
                    page: pl.page
                }).then(function(programs) {
                    return 0 === programs.data.length ? pl.busy = !0 : (_.forEach(programs.data, function(p) {
                        return pl.selected_programs.push(p)
                    }),
                    updateCopyright(),
                    fixProgramBoxHeight(),
                    pl.busy = programs.data.length < pl.limit ? !0 : !1)
                }))
            }
            ,
            updateCopyright = function() {
                return $rootScope.footer.copyrights = _.map(_.uniq(_.pluck(pl.selected_programs, "copyright")), function(c) {
                    return "<p>" + c + "</p>"
                }).join(" ")
            }
            ,
            fixProgramBoxHeight = function() {
                return is_mobile ? $timeout(function() {
                    var pbox;
                    return pbox = angular.element(".program-box"),
                    pbox.css("height", 1.28 * parseInt(pbox.outerWidth(), 10) + "px")
                }, 10) : void 0
            }
            ,
            angular.element($window).bind("resize orientationchange", function() {
                return fixProgramBoxHeight()
            }),
            pl.showPrograms()
        }
        ,
        directive = {
            restrict: "E",
            templateUrl: "app/components/program-list/program-list.html",
            replace: !0,
            scope: {},
            controller: ProgramListController,
            controllerAs: "pl",
            bindToController: {
                programs: "="
            }
        }
    })
}
.call(this),
function() {
    angular.module("hibikiWeb").service("Notification", ["$rootScope", "$timeout", function($rootScope, $timeout) {
        var _notification;
        return _notification = {
            add: function(message, notification_class) {
                var t;
                return t = +new Date,
                $rootScope.notifications[t] = {
                    message: message,
                    "class": notification_class
                },
                $timeout(function() {
                    return delete $rootScope.notifications[t]
                }, 1e4)
            },
            success: function(message) {
                return this.add(message, "success")
            },
            info: function(message) {
                return this.add(message, "info")
            },
            alert: function(message) {
                return this.add(message, "alert")
            }
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").directive("news", ["$sce", "$window", function($sce, $window) {
        var NewsDirectiveController, directive;
        return NewsDirectiveController = function() {
            var news;
            switch (news = this,
            news.label = "",
            news.label_class = "",
            news.data.kind) {
            case 1:
            case 2:
            case 3:
                news.label = "お知らせ",
                news.label_class = "notification";
                break;
            case 4:
                news.label = "ラジオ",
                news.label_class = "radio"
            }
        }
        ,
        directive = {
            restrict: "E",
            replace: !0,
            templateUrl: "app/components/news/news.html",
            scope: {},
            controller: NewsDirectiveController,
            controllerAs: "news",
            bindToController: {
                data: "=",
                large: "="
            }
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").directive("sideNavbar", ["platform", function(platform) {
        var NavbarController, directive;
        return NavbarController = function($timeout) {
            var nb;
            nb = this,
            nb.toggleMenu = !1,
            nb.menuTap = function() {
                return $timeout(function() {
                    return nb.toggleMenu = !1
                }, 300)
            }
            ,
            angular.element("#menuLink").on("click", function(e) {
                return e.preventDefault(),
                angular.element("#menuLink").hasClass("active") ? (angular.element("#layout").removeClass("active"),
                angular.element("#layout .view-container").removeClass("active"),
                angular.element("#menu").removeClass("active"),
                angular.element("#menuLink").removeClass("active"),
                $timeout(function() {
                    return angular.element("#menu").css("display", "none")
                }, 300)) : (angular.element("#menu").css("display", "block"),
                $timeout(function() {
                    return angular.element("#layout").addClass("active"),
                    angular.element("#layout .view-container").addClass("active"),
                    angular.element("#menu").addClass("active"),
                    angular.element("#menuLink").addClass("active")
                }, 100))
            }),
            angular.element("body").on("click", "#layout.active #menu a", function(e) {
                return $timeout(function() {
                    return angular.element("#layout").removeClass("active"),
                    angular.element("#layout .view-container").removeClass("active"),
                    angular.element("#menu").removeClass("active").css("display", ""),
                    angular.element("#menuLink").removeClass("active")
                }, 300)
            })
        }
        ,
        directive = {
            restrict: "A",
            templateUrl: platform().mobile ? "app/components/navbar/navbar.sp.html" : "app/components/navbar/navbar.html",
            scope: {},
            controller: NavbarController,
            controllerAs: "nb",
            bindToController: {
                notifications: "=",
                dailyvoice: "="
            }
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").service("Metatags", ["$rootScope", "meta", "$location", function($rootScope, meta, $location) {
        var defaultMeta;
        return defaultMeta = function() {
            return $rootScope.metatags = {
                title: meta.title,
                keywords: meta.keywords.join(","),
                description: meta.description,
                robots: meta.robots ? meta.robots : "noindex, nofollow",
                ogp: {
                    title: meta.title,
                    description: meta.description,
                    site_name: "響 - HiBiKi Radio Station -",
                    url: $location.absUrl(),
                    image: function() {
                        var host, port, protocol;
                        return protocol = $location.protocol(),
                        host = $location.host(),
                        port = $location.port(),
                        port = 80 === port || 443 === port ? "" : ":" + port,
                        protocol + "://" + host + port + "/assets/images/hibiki_ogp_image.png"
                    }()
                },
                twitter_card_url: {
                    iphone: "jphibikiradio://programs",
                    ipad: "jphibikiradio://programs",
                    googleplay: "jphibikiradio://"
                }
            }
        }
        ,
        $rootScope.metatags || defaultMeta(),
        {
            "default": function() {
                return defaultMeta()
            },
            title: function(new_title) {
                return new_title.length > 0 ? $rootScope.metatags.title = meta.title + " | " + new_title : void 0
            },
            description: function(new_description) {
                return new_description.length > 0 ? $rootScope.metatags.description = new_description : void 0
            },
            keywords: function(new_keywords) {
                return _.isArray(new_keywords) ? $rootScope.metatags.keywords = meta.keywords.concat(new_keywords).join(",") : angular.isString(new_keywords) ? (meta.keywords.push(new_keywords),
                $rootScope.metatags.keywords = meta.keywords.join(",")) : void 0
            },
            robots: function(new_robots) {
                return $rootScope.metatags.robots = new_robots
            },
            ogp: function(prop, val) {
                return $rootScope.metatags.ogp[prop] = val
            },
            twcardurl: function(device, val) {
                return $rootScope.metatags.twitter_card_url[device] = val
            }
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").service("logicaLog", ["$http", "$location", "logicalog_settings", "cookieKeys", "$cookies", "env", "$analytics", function($http, $location, logicalog_settings, cookieKeys, $cookies, env, $analytics) {
        var __logicalog;
        return __logicalog = {
            host: function() {
                return logicalog_settings.host
            },
            uuid: function() {
                return $cookies.get(cookieKeys.uuid)
            },
            session_id: function() {
                return $cookies.get(cookieKeys.sessionid)
            },
            send: function(__params) {
                var endpoint, params, protocol;
                return params = angular.extend(__params, {
                    vuid: this.uuid(),
                    dsid: this.session_id(),
                    env: env,
                    token: "U7JkK5Cd6NWi"
                }),
                protocol = $location.protocol(),
                endpoint = protocol + "://" + this.host() + "/track",
                $http.get(endpoint, {
                    params: params
                }).success(function(data, status, headers, config) {
                    return $analytics.eventTrack(params.action, {
                        category: params.category,
                        label: params.label,
                        value: params.value
                    })
                })
            }
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").directive("logicalogOn", ["logicaLog", function(logicaLog) {
        var isProperty, propertyName;
        return isProperty = function(name) {
            return "logicalog" === name.substr(0, 9) && -1 === ["On", "Event", "If", "Properties"].indexOf(name.substr(9))
        }
        ,
        propertyName = function(name) {
            var s;
            return s = name.slice(9),
            "undefined" != typeof s && null !== s && s.length > 0 ? s.substring(0, 1).toLowerCase() + s.substring(1) : s
        }
        ,
        {
            restrict: "A",
            link: function($scope, $element, $attributes) {
                var __params, eventType;
                return __params = {},
                eventType = $attributes.logicalogOn || "click",
                angular.forEach($attributes.$attr, function(attr, name) {
                    return isProperty(name) ? (__params[propertyName(name)] = $attributes[name],
                    $attributes.$observe(name, function(value) {
                        return __params[propertyName(name)] = value
                    })) : void 0
                }),
                angular.element($element[0]).bind(eventType, function($event) {
                    var eventName;
                    return eventName = $attributes.logicalogEvent || "click",
                    !$attributes.logicalogIf || $scope.$eval($attributes.logicalogIf) ? ($attributes.logicalogProperties && (__params = angular.extend(__params, $scope.$eval($attributes.logicalogProperties))),
                    logicaLog.send(__params)) : void 0
                })
            }
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").service("launchApp", ["$sce", "$window", "$timeout", "platform", "ngDialog", "store_url", function($sce, $window, $timeout, platform, ngDialog, store_url) {
        var __launchApp, appUrl, supportedApp;
        return supportedApp = function() {
            var osv, supported, ua, version;
            return supported = !1,
            ua = navigator.userAgent,
            /(iPhone|iPad|Android)/.test(ua) ? (/Android/.test(ua) ? (version = parseFloat(ua.slice(ua.indexOf("Android") + 8)),
            version >= 4.1 && (supported = !0)) : /(iPhone|iPad)/.test(ua) && (ua.match(/iPhone OS (\w+){1,3}/g),
            osv = (RegExp.$1.replace(/_/g, "") + "00").slice(0, 3),
            osv >= 700 && (supported = !0)),
            supported) : null
        }
        ,
        appUrl = function(path, param) {
            var _app_url, a, isChrome, os, osv, pathinfo, scheme, ua;
            return os = platform().os,
            osv = platform().osv,
            ua = navigator.userAgent.toLowerCase(),
            isChrome = /chrome/.test(ua),
            "android" === os ? (isChrome ? (_app_url = "jphibikiradio://com.bushiroad.HibikiRadio",
            a = $("<a>").attr("href", _app_url).get(0),
            pathinfo = a.pathname.split("/"),
            scheme = a.protocol.substr(0, a.protocol.length - 1),
            _app_url = path && param ? "intent://" + path + "?" + param + "#Intent;scheme=" + scheme + ";package=" + pathinfo[2] + ";end" : "intent://#Intent;scheme=" + scheme + ";package=" + pathinfo[2] + ";end;") : _app_url = path && param ? "jphibikiradio://" + path + "?" + param : "jphibikiradio://",
            _app_url) : "ios" === os ? _app_url = path && param ? "jphibikiradio://" + path + "?" + param : "jphibikiradio://" : void 0
        }
        ,
        __launchApp = {
            store_url: store_url[platform().os],
            dialog: function(app_url, dialogCloseCallback) {
                var _store_url, callback, open_args;
                return _store_url = this.store_url,
                angular.isObject(app_url) && (app_url = appUrl(app_url.path, app_url.params)),
                callback = !1,
                angular.isFunction(dialogCloseCallback) && (callback = !0),
                open_args = {
                    template: "app/components/launch-app/launch-app.html",
                    controller: function($scope) {
                        return $scope.appUrl = app_url,
                        $scope.storeUrl = _store_url,
                        $scope.callback = callback
                    }
                },
                callback && (open_args.preCloseCallback = dialogCloseCallback),
                ngDialog.open(open_args)
            },
            launch: function(path, param) {
                var $iframe, _app_url, isChrome, os, osv, self, start, supported, ua;
                return supported = supportedApp(),
                supported === !1 ? ($window.location.href = this.store_url,
                !1) : (os = platform().os,
                osv = platform().osv,
                ua = navigator.userAgent.toLowerCase(),
                isChrome = /chrome/.test(ua),
                _app_url = appUrl(path, param),
                "android" === os ? isChrome ? (window.location.replace(_app_url),
                self = this,
                $timeout(function() {
                    return self.dialog(_app_url)
                }, 500)) : ($iframe = $("<iframe>"),
                $iframe.addClass("hidden"),
                $iframe.get(0).onload = function() {
                    return $iframe.remove(),
                    window.location.href = this.store_url
                }
                ,
                $iframe.attr("src", _app_url),
                $("body").append($iframe)) : "ios" === os ? osv >= 9 ? this.dialog(_app_url) : ($iframe = $("<iframe>"),
                $iframe.addClass("hidden"),
                self = this,
                start = +new Date,
                $timeout(function() {
                    return +new Date - start < 2e3 ? ($iframe.remove(),
                    window.location.href = self.store_url) : void 0
                }, 500),
                $iframe.attr("src", _app_url),
                $("body").append($iframe)) : void 0)
            }
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").directive("mainInstantVoice", ["platform", function(platform) {
        var directive;
        return directive = {
            restrict: "E",
            templateUrl: platform().mobile ? "app/components/instant-voice/main-instant-voice.sp.html" : "app/components/instant-voice/main-instant-voice.html",
            replace: !1,
            scope: {},
            controller: "InstantVoiceController",
            controllerAs: "iv",
            bindToController: {
                voice: "="
            }
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").directive("mainInstantVoiceCarousel", function() {
        var MainInstantVoiceCarouselController, directive;
        return MainInstantVoiceCarouselController = function(VoicePlayer, platform) {
            var isMobile, iv;
            iv = this,
            iv.isMobile = isMobile = platform().mobile,
            angular.element("#main-instant-voice").on("mouseenter", function() {
                return angular.element("#main-instant-voice slick").slick("slickPause")
            }).on("mouseleave", function() {
                return !VoicePlayer.loading("instant") && VoicePlayer.paused("instant") ? angular.element("#main-instant-voice slick").slick("slickPlay") : void 0
            })
        }
        ,
        directive = {
            restrict: "E",
            templateUrl: "app/components/instant-voice/main-instant-voice-carousel.html",
            replace: !1,
            scope: {},
            controller: MainInstantVoiceCarouselController,
            controllerAs: "iv",
            bindToController: {
                voices: "="
            }
        }
    })
}
.call(this),
function() {
    angular.module("hibikiWeb").directive("instantVoice", function() {
        var directive;
        return directive = {
            restrict: "E",
            templateUrl: "app/components/instant-voice/instant-voice.html",
            replace: !1,
            scope: {},
            controller: "InstantVoiceController",
            controllerAs: "iv",
            bindToController: {
                voice: "=",
                charlength: "=",
                logaction: "="
            }
        }
    })
}
.call(this),
function() {
    angular.module("hibikiWeb").controller("InstantVoiceController", ["$rootScope", "$timeout", "VoicePlayer", "ProgramPlayer", "User", "platform", function($rootScope, $timeout, VoicePlayer, ProgramPlayer, User, platform) {
        var isMobile, isTablet, iv, pauseCarousel, startCarousel;
        iv = this,
        iv.playing = !1,
        iv.loading = !1,
        iv.isMobile = isMobile = platform().mobile,
        iv.isTablet = isTablet = platform().tablet,
        pauseCarousel = function() {
            return angular.element("#main-instant-voice slick").slick("slickPause")
        }
        ,
        startCarousel = function() {
            return angular.element("#main-instant-voice slick").slick("slickPlay")
        }
        ,
        iv.playToggle = function(event) {
            iv.voice.member_limited && User.limited(event) || (event.preventDefault(),
            iv.playing ? (VoicePlayer.pause("instant"),
            iv.loading = !1,
            iv.playing = !1,
            ProgramPlayer.enable() && ProgramPlayer.firstplay && (ProgramPlayer.ctrl = !0)) : iv.voice && iv.voice.video && (VoicePlayer.pauseAll(),
            pauseCarousel(),
            VoicePlayer.voice.instant.loading = !0,
            !isMobile && !isTablet || VoicePlayer.firstplay("instant") || VoicePlayer.load("instant"),
            $timeout(function() {
                var canplay, ended;
                return canplay = function() {
                    return iv.loading = !1,
                    iv.playing = !0,
                    $rootScope.$$phase || $rootScope.$apply()
                }
                ,
                ended = function() {
                    return iv.playing = !1,
                    iv.loading = !1,
                    startCarousel(),
                    $rootScope.$$phase || $rootScope.$apply()
                }
                ,
                iv.loading = !0,
                VoicePlayer.voice.instant.voice.id === iv.voice.id && VoicePlayer.voice.instant.firstplay ? VoicePlayer.playWithFade("instant", canplay, ended) : VoicePlayer.setResource("instant", iv.voice).then(function(url) {
                    return VoicePlayer.playWithFade("instant", canplay, ended)
                })
            }, 0)))
        }
        ,
        iv.voice && iv.voice.video && VoicePlayer.init("instant", iv.voice)
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").filter("urlencode", function() {
        return window.encodeURIComponent
    })
}
.call(this),
function() {
    angular.module("hibikiWeb").filter("truncate", function() {
        return function(text, length, suffix) {
            var b;
            return null == text && (text = ""),
            null == length && (length = 10),
            null == suffix && (suffix = "..."),
            b = 0,
            _.forEach(text, function(s, i) {
                var char_code;
                return char_code = s.charCodeAt(),
                b += char_code >= 65 && 90 >= char_code ? .7 : 255 >= char_code ? .5 : 1,
                b > length ? (text = text.substr(0, i) + suffix,
                !1) : void 0
            }),
            text
        }
    })
}
.call(this),
function() {
    angular.module("hibikiWeb").filter("prefectures_label", ["prefectures", function(prefectures) {
        return function(pref_id) {
            var pref;
            return pref = _.findWhere(prefectures, {
                value: pref_id
            }),
            pref ? pref.label : ""
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").filter("moment", function() {
        return function(dateStr, format) {
            return moment(new Date(dateStr)).format(format)
        }
    }).filter("moment_from_unix", function() {
        return function(timeStr, format) {
            var time;
            return time = parseInt(timeStr, 10),
            moment.unix(time).format(format)
        }
    }).filter("moment_time", function() {
        return function(timeStr) {
            var format, time;
            return time = parseInt(timeStr, 10),
            format = time >= 60 ? "m分 s秒" : "s秒",
            moment.unix(time).format(format)
        }
    })
}
.call(this),
function() {
    angular.module("hibikiWeb").filter("gender_label", ["gender", function(gender) {
        return function(gender_id) {
            var g;
            return g = _.findWhere(gender, {
                value: gender_id
            }),
            g ? g.label : ""
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").filter("allow_mail_label", ["allow_mail", function(allow_mail) {
        return function(value) {
            var a;
            return a = _.findWhere(allow_mail, {
                value: value
            }),
            a ? a.label : ""
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").directive("dailyVoice", ["$rootScope", "$q", "platform", function($rootScope, $q, platform) {
        var DailyVoiceController, directive;
        return DailyVoiceController = function($timeout, $interval, VoicePlayer, ProgramPlayer) {
            var dv, isMobile, isTablet;
            dv = this,
            dv.playing = !1,
            dv.loading = !1,
            dv.isMobile = isMobile = platform().mobile,
            dv.isTablet = isTablet = platform().tablet,
            dv.playToggle = function(event) {
                event.preventDefault(),
                dv.playing ? (VoicePlayer.pause("daily"),
                dv.loading = !1,
                dv.playing = !1,
                ProgramPlayer.enable() && ProgramPlayer.firstplay && (ProgramPlayer.ctrl = !0)) : dv.voice && dv.voice.video && (VoicePlayer.pauseAll(),
                VoicePlayer.voice.daily.player.muted(!1),
                VoicePlayer.voice.daily.loading = !0,
                !isMobile && !isTablet || VoicePlayer.firstplay("daily") || VoicePlayer.load("daily"),
                $timeout(function() {
                    var canplay, ended, playForce;
                    return canplay = function() {
                        return dv.loading = !1,
                        dv.playing = !0,
                        $rootScope.$$phase || $rootScope.$apply()
                    }
                    ,
                    ended = function() {
                        return dv.playing = !1,
                        dv.loading = !1,
                        $rootScope.$$phase || $rootScope.$apply()
                    }
                    ,
                    playForce = function() {
                        return $interval(function() {
                            return VoicePlayer.voice.daily.player.paused() || 0 !== VoicePlayer.voice.daily.player.currentTime() ? void 0 : (VoicePlayer.voice.daily.player.currentTime(1),
                            VoicePlayer.voice.daily.player.currentTime(0))
                        }, 1e3, 10)
                    }
                    ,
                    dv.loading = !0,
                    VoicePlayer.voice.daily.voice.id === dv.voice.id && VoicePlayer.voice.daily.firstplay ? (VoicePlayer.playWithFade("daily", canplay, ended),
                    -1 !== navigator.userAgent.indexOf("Android") ? playForce() : void 0) : VoicePlayer.setResource("daily", dv.voice).then(function(url) {
                        return VoicePlayer.playWithFade("daily", canplay, ended),
                        -1 !== navigator.userAgent.indexOf("Android") ? playForce() : void 0
                    })
                }, 0))
            }
            ,
            dv.voice && dv.voice.video && VoicePlayer.init("daily", dv.voice)
        }
        ,
        directive = {
            restrict: "E",
            templateUrl: platform().mobile ? "app/components/daily-voice/daily-voice.sp.html" : "app/components/daily-voice/daily-voice.html",
            replace: !1,
            scope: {},
            controller: DailyVoiceController,
            controllerAs: "dv",
            bindToController: {
                voice: "="
            }
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").directive("compareTo", function() {
        return {
            require: "ngModel",
            scope: {
                otherModelValue: "=compareTo"
            },
            link: function($scope, $element, $attributes, ngModel) {
                return ngModel.$validators.compareTo = function(modelValue) {
                    return modelValue === $scope.otherModelValue
                }
                ,
                $scope.$watch("otherModelValue", function() {
                    return ngModel.$validate()
                })
            }
        }
    })
}
.call(this),
function() {
    angular.module("hibikiWeb").directive("spBannerCarousel", function() {
        var BannerCarouselSmartphoneController, directive;
        return BannerCarouselSmartphoneController = function($scope, $timeout, $interval, $animate) {
            var bc;
            bc = this,
            bc.$scope = $scope,
            bc.linkTarget = function(b) {
                return angular.isUndefined(b) ? void 0 : 1 === b.kind || 2 === b.kind ? "_blank" : "_self"
            }
        }
        ,
        directive = {
            restrict: "E",
            templateUrl: "app/components/banner-carousel/banner-carousel.sp.html",
            replace: !0,
            scope: {},
            controller: BannerCarouselSmartphoneController,
            controllerAs: "bc",
            bindToController: {
                banners: "="
            }
        }
    })
}
.call(this),
function() {
    angular.module("hibikiWeb").directive("bannerCarousel", function() {
        var BannerCarouselController, directive;
        return BannerCarouselController = function($scope, $timeout, $interval, $animate, $window) {
            var animated, bc, btns, init, intervalId, intervalTime, rotate, setAnimation, thumb_list;
            $window.focus(),
            bc = this,
            bc.$scope = $scope,
            bc.loaded = !1,
            bc.rotate = !1,
            bc.main = [],
            bc.animation = "next",
            bc.animationSpeed = "normal",
            intervalTime = 5e3,
            intervalId = null,
            animated = !1,
            thumb_list = angular.element("#banner-carousel .thumbnail-list"),
            btns = angular.element("#banner-carousel .next-btn, #banner-carousel .prev-btn"),
            init = function() {
                return _.times(2, function() {
                    var last_banner;
                    return last_banner = bc.banners.pop(),
                    bc.banners.unshift(last_banner)
                }),
                setAnimation(),
                bc.next(),
                bc.loaded = !0,
                bc.start(),
                $window.onfocus = function() {
                    return bc.start()
                }
                ,
                $window.onblur = function() {
                    return bc.stop()
                }
            }
            ,
            setAnimation = function() {
                return $timeout(function() {
                    return _.forEach(angular.element("#banner-carousel .thumbnail-list li"), function(b) {
                        return $animate.on("move", b, function(elem, phase) {
                            var main_image, next_main_image;
                            return main_image = angular.element("#banner-carousel .main-banner.current"),
                            next_main_image = angular.element("#banner-carousel .main-banner.next"),
                            "start" === phase ? (animated = !0,
                            btns.hide()) : ($timeout(function() {
                                return thumb_list.addClass("move-active").addClass(bc.animation).addClass(bc.animationSpeed),
                                next_main_image.addClass("current").removeClass("next").addClass(bc.animationSpeed),
                                main_image.addClass("next").removeClass("current").addClass(bc.animationSpeed)
                            }, 0),
                            $timeout(function() {
                                return thumb_list.removeClass("move-next").removeClass("move-prev").removeClass("move-active").removeClass("next").removeClass("prev").removeClass(bc.animation).removeClass(bc.animationSpeed),
                                main_image.removeClass(bc.animationSpeed),
                                next_main_image.removeClass(bc.animationSpeed),
                                animated = !1,
                                btns.show()
                            }, 700))
                        })
                    })
                }, 0)
            }
            ,
            bc.start = function() {
                return bc.stop(),
                intervalId = $interval(function() {
                    return bc.next("normal")
                }, intervalTime)
            }
            ,
            bc.stop = function() {
                return $interval.cancel(intervalId),
                intervalId = null
            }
            ,
            bc.next = function(speed, force) {
                return animated ? void 0 : (null == speed && (speed = "normal"),
                bc.animation = "next",
                bc.animationSpeed = speed,
                rotate())
            }
            ,
            bc.prev = function(speed) {
                return animated ? void 0 : (null == speed && (speed = "normal"),
                bc.animation = "prev",
                bc.animationSpeed = speed,
                rotate())
            }
            ,
            rotate = function() {
                var last_banner, prev_main;
                return prev_main = bc.banners[3],
                "next" === bc.animation ? (last_banner = bc.banners.pop(),
                bc.banners.unshift(last_banner)) : (last_banner = bc.banners.shift(),
                bc.banners.push(last_banner)),
                bc.loaded && thumb_list.addClass("move-" + bc.animation),
                prev_main === bc.main[0] ? bc.main[1] = bc.banners[3] : bc.main[0] = bc.banners[3]
            }
            ,
            bc.$scope.$watch("bc.banners", function(banners) {
                return angular.isUndefined(banners) ? void 0 : init()
            })
        }
        ,
        directive = {
            restrict: "E",
            templateUrl: "app/components/banner-carousel/banner-carousel.html",
            replace: !0,
            scope: {},
            controller: BannerCarouselController,
            controllerAs: "bc",
            bindToController: {
                banners: "="
            }
        }
    })
}
.call(this),
function() {
    angular.module("hibikiWeb").service("Api", ["$http", "apiHost", "serviceStatusAPI", "$q", "platform", function($http, apiHost, serviceStatusAPI, $q, platform) {
        var apiBase, apis, default_paging_options;
        return apiBase = "/api/v1/",
        default_paging_options = function(options) {
            var default_options;
            return null == options && (options = {}),
            default_options = {
                limit: 10,
                page: 1
            },
            angular.extend(default_options, options)
        }
        ,
        apis = {
            service_status: function() {
                return $http.get(serviceStatusAPI).success(function(data, status, headers, config) {})
            },
            news: function(options) {
                var params;
                return params = default_paging_options(options),
                $http.get("" + apiHost + apiBase + "news", {
                    params: params
                }).success(function(data, status, headers, config) {})
            },
            informations: function(options) {
                var params;
                return params = default_paging_options(options),
                $http.get("" + apiHost + apiBase + "informations", {
                    params: params
                }).success(function(data, status, headers, config) {})
            },
            information_detail: function(id) {
                return $http.get("" + apiHost + apiBase + "informations/" + id).success(function(data, status, headers, config) {})
            },
            daily_voice: function() {
                return $http.get("" + apiHost + apiBase + "daily_voice").success(function(data, status, headers, config) {})
            },
            home: function() {
                return $http.get("" + apiHost + apiBase + "home").success(function(data, status, headers, config) {})
            },
            programs: function(options) {
                var params;
                return params = default_paging_options(options),
                "all" === params.day_of_week && delete params.day_of_week,
                $http.get("" + apiHost + apiBase + "/programs", {
                    params: params
                }).success(function(data, status, headers, config) {})
            },
            program_detail: function(access_id) {
                return $http.get("" + apiHost + apiBase + "programs/" + access_id).success(function(data, status, headers, config) {})
            },
            product_informations: function(program_id) {
                return $http.get("" + apiHost + apiBase + "product_informations?program_id=" + program_id).success(function(data, status, headers, config) {})
            },
            program_informations: function(program_id) {
                return $http.get("" + apiHost + apiBase + "program_informations?program_id=" + program_id).success(function(data, status, headers, config) {})
            },
            instant_voices: function(options) {
                var params;
                return params = default_paging_options(options),
                $http.get("" + apiHost + apiBase + "instant_voices", {
                    params: params
                }).success(function(data, status, headers, config) {})
            },
            register: function(data) {
                return $http.post("" + apiHost + apiBase + "users/register", $.param(data)).success(function(data, status, headers, config) {})
            },
            user_profile: function() {
                return $http.get("" + apiHost + apiBase + "/users/profile").success(function(data, status, headers, config) {})
            },
            user_update: function(data) {
                return $http.post("" + apiHost + apiBase + "users/update", $.param(data)).success(function(data, status, headers, config) {})
            },
            password_update: function(data) {
                return $http.post("" + apiHost + apiBase + "users/update_password", $.param(data)).success(function(data, status, headers, config) {})
            },
            withdraw: function() {
                return $http.post("" + apiHost + apiBase + "users/withdraw", $.param({})).success(function(data, status, headers, config) {})
            },
            login: function(data) {
                return $http.post("" + apiHost + apiBase + "users/auth", $.param(data)).success(function(data, status, headers, config) {})
            },
            logout: function() {
                return $http.post("" + apiHost + apiBase + "users/logout", $.param({})).success(function(data, status, headers, config) {})
            },
            password_reminder: function(email) {
                return $http.post("" + apiHost + apiBase + "users/password_reminder", $.param({
                    email: email
                })).success(function(data, status, headers, config) {})
            },
            add_to_favorite: function(program_id) {
                return $http.post("" + apiHost + apiBase + "user_favorites/add", $.param({
                    program_id: program_id
                })).success(function(data, status, headers, config) {})
            },
            delete_to_favorite: function(program_id) {
                return $http.post("" + apiHost + apiBase + "user_favorites/delete", $.param({
                    program_id: program_id
                })).success(function(data, status, headers, config) {})
            },
            favorites: function(options) {
                var params;
                return params = options ? default_paging_options(options) : {},
                $http.get("" + apiHost + apiBase + "user_favorites", {
                    params: params
                }).success(function(data, status, headers, config) {})
            },
            video_url: function(video_id) {
                return $http.get("" + apiHost + apiBase + "videos/play_check?video_id=" + video_id, {
                    cache: !1
                }).then(function(results) {
                    return results.data.playlist_url
                })
            }
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").controller("RegisterController", ["$rootScope", "$scope", "$state", "User", "Notification", "ngDialog", "moment", "logicaLog", function($rootScope, $scope, $state, User, Notification, ngDialog, moment, logicaLog) {
        var _params, _state, reg;
        reg = this,
        $rootScope.savedRoute && $rootScope.savedRoute.name ? (_state = $rootScope.savedRoute.name,
        _params = $rootScope.savedParams) : (_state = "home",
        _params = {}),
        User.loggedIn() && (Notification.info("ログイン済みです"),
        $state.go(_state, _params)),
        reg.prefectures = User.formParts.prefectures,
        reg.birth_date = User.formParts.birth_date,
        reg.user = User.model,
        $scope.$watch(function() {
            return reg.user.birth_date_parts
        }, function() {
            var d, m, y;
            return y = reg.user.birth_date_parts.year,
            m = reg.user.birth_date_parts.month - 1,
            d = reg.user.birth_date_parts.day,
            reg.user.birth_date = moment(new Date(y,m,d)).format("YYYY-MM-DD"),
            User.updateDaysInMonth()
        }, !0),
        reg.showTerms = function(type) {
            return ngDialog.open("member" === type ? {
                template: "app/register/member_terms.html",
                trapFocus: !1,
                preserveFocus: !1
            } : {
                template: "app/register/terms.html",
                trapFocus: !1,
                preserveFocus: !1
            })
        }
        ,
        reg.submit = function() {
            return User.register().then(function(result) {
                return logicaLog.send({
                    category: "complete",
                    action: "register"
                }),
                $state.go("register.complete")
            })["catch"](function(e) {
                var msg;
                return console.log(e),
                msg = e.data.error_message || "新規会員登録に失敗しました。しばらくしてからもう一度お試しください。",
                Notification.alert(msg)
            })
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").controller("ProgramController", ["program_data", "$q", "$state", "platform", "$timeout", "$interval", "$scope", "$rootScope", "$window", "$sce", "ProgramPlayer", "VoicePlayer", "User", "Api", "Notification", "launchApp", "logicaLog", "ngDialog", "moment", function(program_data, $q, $state, platform, $timeout, $interval, $scope, $rootScope, $window, $sce, ProgramPlayer, VoicePlayer, User, Api, Notification, launchApp, logicaLog, ngDialog, moment) {
        var getAdditionalVideo, getVideo, isMobile, isTablet, loadingOn, log_action, program, video_id;
        switch (program = this,
        program.detail = program_data.detail,
        program.products = program_data.products,
        program.informations = program_data.informations,
        program.access_id = $state.params.access_id,
        program.loading = !1,
        program.backstage_loading = !1,
        $rootScope.footer.copyrights = program.detail.copyright,
        $scope.$on("$viewContentLoaded", function(event) {
            return program.detail.copyright ? $rootScope.footer.copyrights = program.detail.copyright : void 0
        }),
        getVideo = function() {
            var _video;
            return _video = null,
            program.detail.episode && program.detail.episode.video && (_video = program.detail.episode.video),
            _video
        }
        ,
        getAdditionalVideo = function() {
            var _video;
            return _video = null,
            program.detail.episode && program.detail.episode.additional_video && (_video = program.detail.episode.additional_video),
            _video
        }
        ,
        program.video = getVideo(),
        program.additional_video = getAdditionalVideo(),
        program.detail.episode && program.detail.episode.html_description && (program.detail.episode.html_description = $sce.trustAsHtml(program.detail.episode.html_description)),
        program.detail.segments && _.each(program.detail.segments, function(s) {
            return s.html_description ? s.html_description = $sce.trustAsHtml(s.html_description) : void 0
        }),
        program.informations && _.each(program.informations, function(info) {
            return info.html_description ? info.html_description = $sce.trustAsHtml(info.html_description) : void 0
        }),
        program.products && _.each(program.products, function(product) {
            return product.html_description ? product.html_description = $sce.trustAsHtml(product.html_description) : void 0
        }),
        program.twitter_hashtags = function() {
            var hash_tag;
            return hash_tag = program.detail.hash_tag,
            0 === hash_tag.length ? "" : _.map(hash_tag.split(" "), function(x) {
                return x.match(/^#/) ? x.slice(1) : x
            }).join(",")
        }(),
        program.isMobile = isMobile = platform().mobile,
        program.isTablet = isTablet = platform().tablet,
        "program" === $state.current.name && $state.go("program.detail", {
            access_id: program.access_id
        }),
        log_action = null,
        $state.current.name) {
        case "program.detail":
            log_action = "view_program";
            break;
        case "program.introductions":
            log_action = "view_program_segments";
            break;
        case "program.products":
            log_action = "view_program_products";
            break;
        case "program.informations":
            log_action = "view_program_informations"
        }
        log_action && logicaLog.send({
            category: "view",
            action: log_action,
            label: program.detail.id
        }),
        program.isSafari = function() {
            var browserText, isSafari, userAgent;
            return isSafari = !1,
            platform().mobile ? isSafari = !1 : (userAgent = window.navigator.userAgent.toLowerCase(),
            browserText = "",
            browserText = -1 !== userAgent.indexOf("msie") || -1 !== userAgent.indexOf("trident") ? "ie" : -1 !== userAgent.indexOf("edge") ? "edge" : -1 !== userAgent.indexOf("chrome") ? "chrome" : -1 !== userAgent.indexOf("safari") ? "safari" : -1 !== userAgent.indexOf("firefox") ? "firefox" : -1 !== userAgent.indexOf("opera") ? "opera" : "",
            "safari" === browserText && (isSafari = !0)),
            isSafari
        }
        ,
        program.stateIs = function(_state) {
            switch (_state) {
            case "detail":
                return "program.detail" === $state.current.name;
            case "introductions":
                return "program.introductions" === $state.current.name;
            case "products":
                return "program.products" === $state.current.name;
            case "informations":
                return "program.informations" === $state.current.name
            }
        }
        ,
        video_id = function(backstage) {
            var __video_id;
            return null == backstage && (backstage = !1),
            __video_id = backstage ? program.detail.episode.additional_video.id : program.detail.episode.video.id
        }
        ,
        loadingOn = function(backstage) {
            return backstage ? program.backstage_loading = !0 : program.loading = !0
        }
        ,
        program.refresh = function() {
            var defered;
            return defered = $q.defer(),
            Api.program_detail($state.params.access_id).then(function(results) {
                return program.detail = results.data,
                program.video = getVideo(),
                program.additional_video = getAdditionalVideo(),
                defered.resolve()
            }),
            defered.promise
        }
        ,
        program.play = function(event, backstage, dialog) {
            event && event.preventDefault(),
            null == dialog && (dialog = !1),
            null == backstage && (backstage = !1),
            ProgramPlayer.player.muted(!1),
            (isMobile || isTablet) && ProgramPlayer.player && !ProgramPlayer.firstplay && ProgramPlayer.player.load(),
            program.refresh().then(function(_this) {
                return function() {
                    var _video, playForce;
                    if (_video = backstage ? program.additional_video : program.video,
                    _video.live_flg) {
                        if (!_.isNull(_video.delivery_start_at) && moment(new Date).isBefore(moment(new Date(_video.delivery_start_at))))
                            return void ngDialog.open({
                                template: '<p style="white-space: pre-line;">' + moment(new Date(_video.delivery_start_at)).format("YYYY年M月D日 H時mm分") + "から配信します。\n配信開始までしばらくお待ち下さい。</p>",
                                plain: !0
                            });
                        if (!_.isNull(_video.delivery_end_at) && moment(new Date).isAfter(moment(new Date(_video.delivery_end_at))))
                            return void ngDialog.open({
                                template: "<p>" + moment(new Date(_video.delivery_end_at)).format("YYYY年M月D日 H時mm分") + "に配信を終了しました</p>",
                                plain: !0
                            })
                    }
                    return isMobile && dialog && !_this.playing() ? (ProgramPlayer.player || ProgramPlayer.init(program.detail, backstage),
                    void launchApp.dialog({
                        path: "program",
                        params: "program_id=" + program.detail.id + "&episode_id=" + program.detail.episode.id + "&target_video_id=" + video_id(backstage)
                    }, function(play) {
                        return "play" === play && program.play(null, backstage, !1),
                        !0
                    })) : backstage && User.limited(event) ? void 0 : (playForce = function() {
                        return $interval(function() {
                            return ProgramPlayer.player.paused() || 0 !== ProgramPlayer.player.currentTime() ? void 0 : (ProgramPlayer.player.currentTime(1),
                            ProgramPlayer.player.currentTime(0))
                        }, 1e3, 10)
                    }
                    ,
                    ProgramPlayer.player ? (VoicePlayer.pauseAll(),
                    ProgramPlayer.program_id === program.detail.id && ProgramPlayer.video_id === video_id(backstage) && ProgramPlayer.firstplay ? ProgramPlayer.playToggle() : (loadingOn(backstage),
                    ProgramPlayer.setResource(program.detail, backstage).then(function() {
                        return ProgramPlayer.play(),
                        -1 !== navigator.userAgent.indexOf("Android") ? playForce() : void 0
                    }))) : (loadingOn(backstage),
                    ProgramPlayer.init(program.detail, backstage).then(function() {
                        return ProgramPlayer.play(),
                        -1 !== navigator.userAgent.indexOf("Android") ? playForce() : void 0
                    })))
                }
            }(this))["catch"](function(e) {
                return $state.go("error")
            })
        }
        ,
        program.playing = function(backstage) {
            return null == backstage && (backstage = !1),
            ProgramPlayer.program_id !== program.detail.id || null === ProgramPlayer.player || ProgramPlayer.player.paused() ? !1 : backstage && ProgramPlayer.video_id === video_id(backstage) ? (program.backstage_loading = !1,
            !0) : backstage || ProgramPlayer.video_id !== video_id(backstage) ? !1 : (program.loading = !1,
            !0)
        }
        ,
        program.letter = function(url) {
            var height, options, width, wsh, wsw;
            width = 840,
            height = 750,
            options = "width=" + width + ",height=" + height + ",menubar=no,toolbar=no,scrollbars=yes,status=no,location=no",
            wsw = $window.screen.availWidth,
            wsh = $window.screen.availHeight,
            "undefined" != typeof $window.screenX ? options += ",screenX=" + ($window.screenX = (wsw - width) / 2) + ",screenY=" + ($window.screenY + (wsh - height) / 2) : "undefined" != typeof $window.screenLeft && (options += ",left=" + ($window.screenLeft + (wsw - width) / 2) + ",top=" + ($window.screenTop + (wssh - height) / 2)),
            $window.open(url, "letterWindow", options)
        }
        ,
        program.addToFavorite = function(event) {
            return User.limited(event) ? void 0 : Api.add_to_favorite(program.detail.id).then(function(results) {
                return program.detail.user_favorite_flg = !0,
                Notification.success("お気に入りに追加しました")
            })["catch"](function(e) {
                var msg;
                return msg = e.data.error_message || "エラーが発生しました",
                Notification.alert(msg)
            })
        }
        ,
        program.deleteToFavorite = function(event) {
            return User.limited(event) ? void 0 : Api.delete_to_favorite(program.detail.id).then(function(results) {
                return program.detail.user_favorite_flg = !1,
                Notification.success("お気に入りから削除しました")
            })["catch"](function(e) {
                var msg;
                return msg = e.data.error_message || "エラーが発生しました",
                Notification.alert(msg)
            })
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").controller("PasswordReminderController", ["User", "Notification", "$state", function(User, Notification, $state) {
        var pr;
        pr = this,
        pr.email = null,
        pr.submit = function() {
            return User.password_reminder(pr.email).then(function(result) {
                return console.log(result),
                Notification.success("入力したメールアドレス宛にパスワード変更手続きメールを送信しました。"),
                $state.go("login")
            })["catch"](function(e) {
                return Notification.alert(400 === e.status && e.data.result && 22 === e.data.result ? e.data.error_message : "エラーが発生しました")
            })
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").controller("NotificationsController", ["data", "Api", function(data, Api) {
        var nt;
        nt = this,
        nt.busy = !1,
        nt.page = 1,
        nt.limit = 10,
        nt.notifications = data.data,
        nt.loadMore = function() {
            return nt.busy ? void 0 : (nt.busy = !0,
            nt.page += 1,
            Api.news({
                limit: nt.limit,
                page: nt.page
            }).then(function(more_data) {
                return 0 === more_data.data.length ? nt.busy = !0 : (_.forEach(more_data.data, function(n) {
                    return nt.notifications.push(n)
                }),
                nt.busy = !1)
            }))
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").controller("MypageController", ["$scope", "$state", "$window", "$timeout", "User", "Notification", "ngDialog", "logicaLog", function($scope, $state, $window, $timeout, User, Notification, ngDialog, logicaLog) {
        var _countdown, my;
        my = this,
        my.regist_user = User.get(),
        my.prefectures = User.formParts.prefectures,
        my.birth_date = User.formParts.birth_date,
        my.user = User.model,
        my.user.email = my.regist_user.email,
        my.user.name = my.regist_user.name,
        my.user.birth_date = my.regist_user.birth_date,
        my.user.birth_date_parts = {
            year: moment(new Date(my.user.birth_date)).year(),
            month: moment(new Date(my.user.birth_date)).month() + 1,
            day: moment(new Date(my.user.birth_date)).date()
        },
        my.user.prefecture_id = my.regist_user.prefecture_id,
        my.user.allow_mail_flg = my.regist_user.allow_mail_flg ? "true" : "false",
        my.user.gender = my.regist_user.gender,
        $scope.$watch(function() {
            return my.user.birth_date_parts
        }, function() {
            var d, m, y;
            return y = my.user.birth_date_parts.year,
            m = my.user.birth_date_parts.month - 1,
            d = my.user.birth_date_parts.day,
            my.user.birth_date = moment(new Date(y,m,d)).format("YYYY-MM-DD"),
            User.updateDaysInMonth()
        }, !0),
        my.editSubmit = function() {
            User.update().then(function(result) {
                return Notification.success("会員情報を更新しました"),
                $state.go("mypage")
            })["catch"](function(e) {
                var msg;
                return console.log(e),
                msg = e.data.error_message || "エラーが発生しました",
                Notification.alert(msg)
            })
        }
        ,
        my.password = {
            old: "",
            "new": "",
            new_confirm: ""
        },
        my.passwordChangeSubmit = function() {
            User.password_update(my.password.old, my.password["new"]).then(function(result) {
                return Notification.success("パスワードを変更しました"),
                logicaLog.send({
                    category: "complete",
                    action: "password_change"
                }),
                $state.go("mypage")
            })["catch"](function(e) {
                var msg;
                return console.log(e),
                msg = e.data.error_message || "エラーが発生しました",
                Notification.alert(msg)
            })
        }
        ,
        my.withdraw_exec = !1,
        my.countdown = 10,
        _countdown = function() {
            return $timeout(function() {
                return my.countdown -= 1,
                my.countdown > 0 ? _countdown() : void 0
            }, 1e3)
        }
        ,
        my.withdrawConfirm = function() {
            ngDialog.open({
                template: "app/mypage/partials/_withdraw_confirm.html",
                trapFocus: !1,
                preserveFocus: !1,
                controller: function($scope) {
                    return $scope.withdraw = function() {
                        return my.withdraw(),
                        $scope.closeThisDialog()
                    }
                }
            })
        }
        ,
        my.withdraw = function() {
            User.withdraw().then(function() {
                return my.withdraw_exec = !0,
                logicaLog.send({
                    category: "complete",
                    action: "withdraw"
                }),
                _countdown(),
                $timeout(function() {
                    return $window.location.href = "/"
                }, 10500)
            })["catch"](function(e) {
                var msg;
                return console.log(e),
                msg = e.data.error_message || "エラーが発生しました",
                Notification.alert(msg)
            })
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").controller("MainController", ["$rootScope", "$scope", "$timeout", "data", "notifications", function($rootScope, $scope, $timeout, data, notifications) {
        var main;
        main = this,
        main.banners = data.data.top_banners,
        main.voice = _.first(data.data.instant_voices),
        main.voices = data.data.instant_voices,
        main.news = notifications.data
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").controller("LogoutController", ["$rootScope", "$state", "$q", "$timeout", "$window", "User", "Notification", function($rootScope, $state, $q, $timeout, $window, User, Notification) {
        var _countdown, logout;
        logout = this,
        logout.countdown = 3,
        _countdown = function() {
            return $timeout(function() {
                return logout.countdown -= 1,
                logout.countdown > 0 ? _countdown() : void 0
            }, 1e3)
        }
        ,
        User.logout().then(function(result) {
            return Notification.success("ログアウトしました"),
            _countdown(),
            $timeout(function() {
                return $window.location.href = "/"
            }, 3500)
        })["catch"](function(e) {
            return 401 === e.status ? (Notification.info(e.data.error_message),
            User.destroy(),
            $state.go($rootScope.previousState)) : (Notification.alert("エラーが発生しました"),
            $state.go("home"))
        })
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").controller("LoginController", ["$q", "$timeout", "$rootScope", "$state", "User", "Api", "ProgramPlayer", "Notification", function($q, $timeout, $rootScope, $state, User, Api, ProgramPlayer, Notification) {
        var _params, _state, login;
        login = this,
        $rootScope.savedRoute && $rootScope.savedRoute.name ? (_state = $rootScope.savedRoute.name,
        _params = $rootScope.savedParams) : (_state = "home",
        _params = {}),
        User.loggedIn() && (Notification.info("ログイン済みです"),
        $state.go(_state, _params)),
        login.submit = function() {
            User.login({
                email: this.email,
                password: this.password
            }).then(function(result) {
                var backstage, before_time, deferred, is_paused;
                return deferred = $q.defer(),
                ProgramPlayer.player && ProgramPlayer.program_id ? (backstage = ProgramPlayer.backstage,
                is_paused = ProgramPlayer.paused(),
                before_time = ProgramPlayer.player.currentTime(),
                Api.program_detail(ProgramPlayer.program_id).then(function(result) {
                    return ProgramPlayer.setResource(result.data, backstage).then(function() {
                        return $timeout(function() {
                            return before_time > 0 && ProgramPlayer.player.currentTime(before_time),
                            is_paused || ProgramPlayer.play(),
                            deferred.resolve()
                        }, 500)
                    })
                })) : deferred.resolve(),
                deferred.promise
            }).then(function(result) {
                return Notification.success("ログインしました"),
                $state.go(_state, _params)
            })["catch"](function(e) {
                return Notification.alert(401 === e.status ? "メールアドレスまたはパスワードが違います" : e.data.error_message)
            })
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").controller("InstantVoicesController", ["data", "Api", function(data, Api) {
        var iv;
        iv = this,
        iv.busy = !1,
        iv.page = 1,
        iv.limit = 12,
        iv.instant_voices = data.data,
        iv.loadMore = function() {
            return iv.busy ? void 0 : (iv.busy = !0,
            iv.page += 1,
            Api.instant_voices({
                limit: iv.limit,
                page: iv.page
            }).then(function(voices) {
                return 0 === voices.data.length ? iv.busy = !0 : (_.forEach(voices.data, function(v) {
                    return iv.instant_voices.push(v)
                }),
                iv.busy = !1)
            }))
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").controller("InformationsController", ["data", "Api", function(data, Api) {
        var info;
        info = this,
        info.busy = !1,
        info.page = 1,
        info.limit = 10,
        info.informations = data.data,
        info.labelClass = function(kind) {
            var label_class;
            switch (kind) {
            case 1:
                label_class = "program";
                break;
            case 2:
                label_class = "event";
                break;
            case 3:
                label_class = "article";
                break;
            case 4:
                label_class = "other"
            }
            return label_class
        }
        ,
        info.kinds = [{
            value: "",
            label: "全て"
        }, {
            value: "1",
            label: "番組情報"
        }, {
            value: "2",
            label: "イベント情報"
        }, {
            value: "3",
            label: "商品情報"
        }, {
            value: "4",
            label: "その他お知らせ"
        }],
        info.selectedKind = "",
        info.searchKind = function() {
            return info.page = 1,
            Api.informations({
                limit: info.limit,
                page: info.page,
                kind: info.selectedKind
            }).then(function(informations) {
                return info.informations = informations.data
            })
        }
        ,
        info.loadMore = function() {
            return info.busy ? void 0 : (info.busy = !0,
            info.page += 1,
            Api.informations({
                limit: info.limit,
                page: info.page,
                kind: info.selectedKind
            }).then(function(informations) {
                return 0 === informations.data.length ? info.busy = !0 : (_.forEach(informations.data, function(i) {
                    return info.informations.push(i)
                }),
                info.busy = !1)
            }))
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").controller("InformationDetailController", ["$sce", "data", "platform", "Metatags", function($sce, data, platform, Metatags) {
        var info_d;
        info_d = this,
        info_d.data = data.data,
        info_d.isMobile = platform().mobile,
        Metatags.title(info_d.data.name),
        info_d.data && info_d.data.html_description && (info_d.data.html_description = $sce.trustAsHtml(info_d.data.html_description))
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").controller("HeaderController", ["notifications", "daily_voice", function(notifications, daily_voice) {
        var header;
        header = this,
        header.notifications = notifications.data,
        header.daily_voice = daily_voice ? daily_voice.data : null
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").controller("FooterController", ["$rootScope", "$scope", "store_url", function($rootScope, $scope, store_url) {
        var footer;
        null == $rootScope.footer && ($rootScope.footer = {}),
        footer = this,
        footer.store_url = store_url,
        $scope.$watch(function() {
            return $rootScope.footer.copyrights
        }, function() {
            return footer.copyright = $rootScope.footer.copyrights
        })
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").controller("FavoriteController", ["$rootScope", "$scope", "data", "Api", "User", "Notification", "ProgramQueue", "ProgramPlayer", "platform", "$timeout", function($rootScope, $scope, data, Api, User, Notification, ProgramQueue, ProgramPlayer, platform, $timeout) {
        var fav, updateCopyright;
        fav = this,
        fav.busy = !1,
        fav.page = 1,
        fav.limit = 10,
        fav.editable = !1,
        fav.program = [],
        updateCopyright = function() {
            return $rootScope.footer.copyrights = _.uniq(_.pluck(fav.programs, "copyright")).join(" ")
        }
        ,
        $scope.$watch("fav.programs", function() {
            return updateCopyright()
        }),
        fav.programs = _.pluck(data.data, "program"),
        fav.playWithQueue = function(event, program_id) {
            return event.preventDefault(),
            Api.favorites(!1).then(function(result) {
                var _favorite_programs, favorite_programs, tmp_list;
                return _favorite_programs = _.pluck(result.data, "program"),
                tmp_list = [],
                _.forEach(_favorite_programs, function(f) {
                    return f.id === program_id ? !1 : tmp_list.push(f)
                }),
                favorite_programs = _.difference(_favorite_programs, tmp_list).concat(tmp_list),
                ProgramQueue.set(favorite_programs),
                ProgramQueue.play()
            })
        }
        ,
        fav.play = function(event, access_id) {
            var isMobile;
            return event.preventDefault(),
            isMobile = platform().mobile,
            isMobile && !ProgramPlayer.firstplay && ProgramPlayer.player.load(),
            $timeout(function() {
                return Api.program_detail(access_id).then(function(result) {
                    var program;
                    return program = result.data,
                    program.episode && program.episode.video ? ProgramPlayer.setResource(program).then(function() {
                        return ProgramPlayer.play()
                    }) : Notification.info("現在再生可能なエピソードはありません")
                })
            }, 0)
        }
        ,
        fav.deleteFavorite = function(event, program_id) {
            return User.limited(event) ? void 0 : Api.delete_to_favorite(program_id).then(function(results) {
                return _.remove(fav.programs, function(n) {
                    return n.id === program_id
                }),
                Notification.success("お気に入りから削除しました")
            })["catch"](function(e) {
                return Notification.alert("エラーが発生しました")
            })
        }
        ,
        fav.loadMore = function() {
            fav.busy || (fav.busy = !0,
            fav.page += 1,
            Api.favorites({
                limit: fav.limit,
                page: fav.page
            }).then(function(favorites) {
                return 0 === favorites.data.length ? fav.busy = !0 : (_.forEach(favorites.data, function(i) {
                    return fav.programs.push(i.program)
                }),
                fav.busy = !1)
            }))
        }
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").run(["$log", "$state", "$rootScope", "$http", "$window", "$location", "$timeout", "$resource", "CacheFactory", "platform", "hasFlash", "Metatags", "User", "Notification", "cookieKeys", "$cookies", "uuid4", "logicaLog", "env", "launchApp", "ProgramPlayer", "VoicePlayer", "Api", function($log, $state, $rootScope, $http, $window, $location, $timeout, $resource, CacheFactory, platform, hasFlash, Metatags, User, Notification, cookieKeys, $cookies, uuid4, logicaLog, env, launchApp, ProgramPlayer, VoicePlayer, Api) {
        var _platform, _sid, _uuid, externalLintTracking, scroll_adjust;
        return $log.debug("runBlock start"),
        _uuid = $cookies.get(cookieKeys.uuid),
        _uuid || $cookies.put(cookieKeys.uuid, uuid4.generate(), {
            expires: moment(new Date).add(10, "years").toDate()
        }),
        _sid = $cookies.get(cookieKeys.sessionid),
        _sid || $cookies.put(cookieKeys.sessionid, uuid4.generate()),
        $rootScope.env = env,
        $rootScope.hasFlash = hasFlash(),
        null == $rootScope.footer && ($rootScope.footer = {}),
        $rootScope.footer.copyrights = "",
        $rootScope.User = User,
        null == $rootScope.notifications && ($rootScope.notifications = {}),
        $rootScope.htmlReady = function() {
            return $rootScope.$evalAsync(function() {
                return $timeout(function() {
                    return angular.isFunction(window.callPhantom) ? window.callPhantom() : void 0
                }, 0)
            })
        }
        ,
        $rootScope.platform = _platform = platform(),
        _platform.mobile ? angular.element("body").addClass("mobile").addClass(_platform.os) : (angular.element("body").addClass("pc"),
        _platform.tablet && angular.element("body").addClass("tablet").addClass(_platform.os)),
        _platform.mobile || (scroll_adjust = function(ev) {
            var $content, $header, $sidebar, $target, _width, bottom, content_height, controller_height, height_diff, scrollTop, sidebar_height, window_height;
            $header = angular.element("body.pc #layout header"),
            $sidebar = angular.element("body.pc #layout header #menu"),
            $content = angular.element("body.pc #layout .view-container"),
            $sidebar && $content && ($target = angular.element(ev.currentTarget),
            sidebar_height = $sidebar.height(),
            content_height = $content.height(),
            scrollTop = $target.scrollTop(),
            controller_height = angular.element("#player-ctrl.active").height(),
            null === controller_height && (controller_height = 0),
            window_height = $target.height() - controller_height,
            content_height > sidebar_height ? ($content.removeClass("fix").css({
                width: "",
                bottom: ""
            }),
            window_height + scrollTop >= sidebar_height ? $header.addClass("fix").css({
                bottom: $sidebar.height() + "px"
            }) : $header.removeClass("fix").css({
                bottom: ""
            })) : ($header.removeClass("fix").css({
                bottom: ""
            }),
            window_height + scrollTop >= content_height ? (_width = $content.width(),
            height_diff = window_height - $content.height(),
            bottom = 0 > height_diff ? 0 : height_diff,
            $content.addClass("fix").css({
                width: _width,
                bottom: bottom + "px"
            })) : $content.removeClass("fix").css({
                width: "",
                bottom: ""
            })))
        }
        ,
        _platform.tablet ? "ios" === _platform.os && _platform.osv <= 7 ? angular.element("#layout").on("scroll resize touchmove", _.throttle(scroll_adjust, 200)) : angular.element("#layout").on("scroll resize", scroll_adjust) : angular.element($window).on("scroll resize", scroll_adjust)),
        externalLintTracking = function(action) {
            return function() {
                var url;
                return url = angular.element(this).attr("href"),
                logicaLog.send({
                    category: "tap_link",
                    action: action,
                    label: url
                })
            }
        }
        ,
        angular.element("body").on("click", "#program-detail a[target='_blank']", externalLintTracking("program_detail")),
        angular.element("body").on("click", "#program-introductions a[target='_blank']", externalLintTracking("program_segments")),
        angular.element("body").on("click", "#program-products a[target='_blank']", externalLintTracking("product_information")),
        angular.element("body").on("click", "#program-informations a[target='_blank']", externalLintTracking("program_information")),
        angular.element("body").on("click", "#information-detail a[target='_blank']", externalLintTracking("information")),
        $rootScope.$on("$stateNotFound", function(event, unfoundState, fromState, fromParams) {
            return event.preventDefault(),
            Notification.alert(unfoundState.to + " is undefined!")
        }),
        $rootScope.$on("$stateChangeError", function(event, toState, toParams, fromState, fromParams, error) {
            return event.preventDefault(),
            403 === error.status ? (Notification.info("ログインが必要です"),
            User.destroy(),
            $rootScope.savedRoute = toState,
            $rootScope.savedParams = toParams,
            $state.go("login")) : 500 === error.status ? Notification.alert("エラーが発生しました") : void 0
        }),
        $rootScope.$on("$stateChangeStart", function(event, toState, toParams, fromState, fromParams) {
            var _query;
            return "maintenance" !== toState.name && Api.service_status().then(function(result) {
                return $rootScope.terms_url = result.data.terms_url,
                $rootScope.member_terms_url = result.data.member_terms_url,
                $rootScope.privacy_policy_url = result.data.privacy_policy_url,
                2 === result.data.status ? ($rootScope.maintenance_reason = result.data.maintenance.reason,
                $rootScope.maintenance_url = result.data.maintenance.page_url,
                $state.go("maintenance")) : void 0
            }),
            $resource("./build.json").get().$promise.then(function(res) {
                var build_version, local_build_version;
                return build_version = parseInt(res.build_version, 10),
                local_build_version = parseInt($cookies.get("build_version"), 10),
                ProgramPlayer.paused() && build_version !== local_build_version && ($cookies.put("build_version", build_version),
                !isNaN(local_build_version)) ? $window.location.href = $state.href(toState.name, toParams) : void 0
            }),
            $rootScope.standalone = !1,
            _query = $location.search(),
            _query.standalone && ($rootScope.standalone = !0),
            VoicePlayer.voice.instant && VoicePlayer.voice.instant.pause(),
            toState.auth ? User.loggedIn() ? void 0 : (event.preventDefault(),
            Notification.info("ログインが必要です"),
            $rootScope.savedRoute = toState,
            $rootScope.savedParams = toParams,
            $state.go("login")) : void 0
        }),
        $rootScope.$on("$stateChangeSuccess", function(event, toState, toParams, fromState, fromParams, error) {
            var params;
            return delete window.google_persistent_state_async,
            "^" === fromState.url && ($cookies.remove(cookieKeys.sessinid),
            $cookies.put(cookieKeys.sessionid, uuid4.generate())),
            toState.log_action && (params = {
                category: "view",
                action: toState.log_action
            },
            toState.log_label && toParams[toState.log_label] && (params.label = toParams[toState.log_label]),
            logicaLog.send(params)),
            $rootScope.currentUrl = $location.absUrl(),
            toState.parent && "app" === toState.parent && (Metatags["default"](),
            toState.title && Metatags.title(toState.title),
            toState.title && Metatags.ogp("title", toState.title)),
            Metatags.ogp("url", $location.absUrl()),
            /^program/.test(toState.name) && !/^program/.test(fromState.name) && angular.element("html, body").scrollTop(0),
            (/^login/.test(toState.name) || /^register/.test(toState.name)) && fromState.name.length > 0 && ($rootScope.savedRoute = fromState,
            $rootScope.savedParams = fromParams),
            $rootScope.previousState = fromState.name,
            $rootScope.currentState = toState.name,
            $rootScope.footer.copyrights = '<i class="fa fa-copyright"></i> bushiroad All Rights Reserved.',
            "undefined" != typeof twttr ? $timeout(function() {
                return twttr.widgets.load()
            }, 0) : void 0
        }),
        angular.module("infinite-scroll").value("THROTTLE_MILLISECONDS", 250),
        angular.element("body").on("contextmenu", function() {
            return !1
        }),
        $log.debug("runBlock end")
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").config(["$stateProvider", "$urlRouterProvider", "platform", function($stateProvider, $urlRouterProvider, platform) {
        var _platform;
        return _platform = platform(),
        $stateProvider.state("app", {
            "abstract": !0,
            views: {
                "@": {
                    template: "<ui-view />"
                },
                "header@": {
                    templateUrl: _platform.mobile ? "app/header/header.sp.html" : "app/header/header.html",
                    controller: "HeaderController",
                    controllerAs: "header"
                },
                "footer@": {
                    templateUrl: _platform.mobile ? "app/footer/footer.sp.html" : "app/footer/footer.html",
                    controller: "FooterController",
                    controllerAs: "footer"
                }
            },
            resolve: {
                notifications: ["Api", function(Api) {
                    return Api.news({
                        limit: 5
                    })
                }
                ],
                daily_voice: ["Api", function(Api) {
                    return Api.daily_voice().then(function(result) {
                        return result
                    })["catch"](function(e) {
                        return 400 === e.status ? null : void 0
                    })
                }
                ]
            }
        }).state("home", {
            parent: "app",
            title: "響 - HiBiKi Radio Station -",
            url: "/",
            log_action: "view_home",
            views: {
                "@": {
                    templateUrl: _platform.mobile ? "app/main/main.sp.html" : "app/main/main.html",
                    controller: "MainController",
                    controllerAs: "main"
                }
            },
            resolve: {
                data: ["Api", function(Api) {
                    return Api.home()
                }
                ]
            }
        }).state("program", {
            parent: "app",
            url: "/description/:access_id",
            views: {
                "@": {
                    templateUrl: _platform.mobile ? "app/program/program.sp.html" : "app/program/program.html",
                    controller: "ProgramController",
                    controllerAs: "program"
                }
            },
            resolve: {
                program_data: ["$timeout", "$state", "$stateParams", "$location", "Api", "platform", "ProgramPlayer", "moment", "Metatags", "launchApp", function($timeout, $state, $stateParams, $location, Api, platform, ProgramPlayer, moment, Metatags, launchApp) {
                    var access_id, program;
                    return access_id = $stateParams.access_id,
                    program = {},
                    Api.program_detail(access_id).then(function(results) {
                        return program.detail = results.data,
                        Api.product_informations(program.detail.id)
                    }).then(function(results) {
                        return program.products = results.data,
                        Api.program_informations(program.detail.id)
                    }).then(function(results) {
                        return program.informations = results.data,
                        program.loaded = !0
                    }).then(function() {
                        var _video, cast_names, keywords, title;
                        return platform().mobile && $location.hash().match(/^share/) && launchApp.launch("program", "program_id=" + program.detail.id),
                        ProgramPlayer.player || program.detail.episode && program.detail.episode.video && (_video = program.detail.episode.video,
                        (!_video.live_flg || _video.live_flg && moment(new Date).isBetween(moment(new Date(_video.delivery_start_at)), moment(new Date(_video.delivery_end_at)))) && ProgramPlayer.init(program.detail)),
                        title = program.detail.meta_title.length > 0 ? program.detail.meta_title : program.detail.name,
                        cast_names = _.pluck(program.detail.casts, "name"),
                        keywords = cast_names.concat(program.detail.meta_keyword.split(",")),
                        Metatags.title(title),
                        Metatags.keywords(keywords),
                        Metatags.description(program.detail.meta_description),
                        Metatags.ogp("title", title),
                        Metatags.ogp("description", program.detail.meta_description || program.detail.description),
                        Metatags.ogp("image", program.detail.pc_image_url),
                        Metatags.twcardurl("iphone", "jphibikiradio://programs?program_id=" + program.detail.id),
                        Metatags.twcardurl("ipad", "jphibikiradio://programs?program_id=" + program.detail.id),
                        Metatags.twcardurl("googleplay", "jphibikiradio://programs?program_id=" + program.detail.id),
                        program
                    })["catch"](function(e) {
                        return $state.go("error")
                    })
                }
                ]
            }
        }).state("program.detail", {
            parent: "program",
            url: "/detail",
            templateUrl: "app/program/partials/_detail.html"
        }).state("program.introductions", {
            parent: "program",
            url: "/introductions",
            templateUrl: "app/program/partials/_introductions.html"
        }).state("program.informations", {
            parent: "program",
            url: "/informations",
            templateUrl: "app/program/partials/_informations.html"
        }).state("program.products", {
            parent: "program",
            url: "/products",
            templateUrl: "app/program/partials/_products.html"
        }).state("notifications", {
            parent: "app",
            title: "最新情報",
            url: "/notifications",
            log_action: "view_notifications",
            views: {
                "@": {
                    templateUrl: _platform.mobile ? "app/notifications/notifications.sp.html" : "app/notifications/notifications.html",
                    controller: "NotificationsController",
                    controllerAs: "nt"
                }
            },
            resolve: {
                data: ["Api", function(Api) {
                    return Api.news({
                        limit: 10
                    })
                }
                ]
            }
        }).state("informations", {
            parent: "app",
            title: "お知らせ一覧",
            url: "/informations",
            log_action: "view_informations",
            views: {
                "@": {
                    templateUrl: "app/informations/informations.html",
                    controller: "InformationsController",
                    controllerAs: "info"
                }
            },
            resolve: {
                data: ["Api", function(Api) {
                    return Api.informations({
                        limit: 10
                    })
                }
                ]
            }
        }).state("informations.detail", {
            url: "/:info_id",
            log_action: "view_event",
            log_label: "info_id",
            views: {
                "@": {
                    templateUrl: "app/informations/information_detail.html",
                    controller: "InformationDetailController",
                    controllerAs: "info_d"
                }
            },
            resolve: {
                data: ["$stateParams", "Api", function($stateParams, Api) {
                    var info_id;
                    return info_id = $stateParams.info_id,
                    Api.information_detail(info_id)
                }
                ]
            }
        }).state("instantvoices", {
            parent: "app",
            title: "インスタンスボイス一覧",
            url: "/instantvoices",
            log_action: "view_instantvoices",
            views: {
                "@": {
                    templateUrl: "app/instantvoices/instantvoices.html",
                    controller: "InstantVoicesController",
                    controllerAs: "iv"
                }
            },
            resolve: {
                data: ["Api", function(Api) {
                    return Api.instant_voices({
                        limit: 12
                    })
                }
                ]
            }
        }).state("favorites", {
            parent: "app",
            title: "お気に入り一覧",
            auth: !0,
            url: "/favorites",
            log_action: "view_user_favorites",
            views: {
                "@": {
                    templateUrl: _platform.mobile ? "app/favorites/favorites.sp.html" : "app/favorites/favorites.html",
                    controller: "FavoriteController",
                    controllerAs: "fav"
                }
            },
            resolve: {
                data: ["Api", "platform", "ProgramPlayer", function(Api, platform, ProgramPlayer) {
                    return Api.favorites({
                        limit: 10
                    }).then(function(result) {
                        var program;
                        return !ProgramPlayer.player && result.data && result.data.length > 0 && (program = result.data[0].program,
                        program && program.episode && program.episode.video && ProgramPlayer.init(program)),
                        result
                    })
                }
                ]
            }
        }).state("login", {
            parent: "app",
            title: "ログイン",
            url: "/login",
            log_action: "view_login",
            views: {
                "@": {
                    templateUrl: "app/login/login.html",
                    controller: "LoginController",
                    controllerAs: "login"
                }
            }
        }).state("logout", {
            parent: "app",
            title: "ログアウト",
            auth: !0,
            url: "/logout",
            log_action: "view_logout",
            views: {
                "@": {
                    templateUrl: "app/logout/logout.html",
                    controller: "LogoutController",
                    controllerAs: "logout"
                }
            }
        }).state("password_reminder", {
            parent: "app",
            title: "パスワードリマインダー",
            url: "/password_reminder",
            log_action: "view_password_reminder",
            views: {
                "@": {
                    templateUrl: "app/password_reminder/password_reminder.html",
                    controller: "PasswordReminderController",
                    controllerAs: "pr"
                }
            }
        }).state("mypage", {
            parent: "app",
            auth: !0,
            url: "/mypage",
            log_action: "view_mypage",
            views: {
                "@": {
                    templateUrl: "app/mypage/mypage.html",
                    controller: "MypageController",
                    controllerAs: "my"
                }
            },
            resolve: {
                profile: ["User", function(User) {
                    return User.profile()
                }
                ]
            }
        }).state("mypage.profile_edit", {
            auth: !0,
            url: "/profile_edit",
            log_action: "view_profile_edit",
            templateUrl: "app/mypage/partials/_profile_edit.html"
        }).state("mypage.password", {
            auth: !0,
            url: "/password_change",
            log_action: "view_password_change",
            templateUrl: "app/mypage/partials/_password_change.html"
        }).state("mypage.withdraw", {
            auth: !0,
            url: "/withdraw",
            log_action: "view_withdraw",
            templateUrl: "app/mypage/partials/_withdraw.html"
        }).state("register", {
            parent: "app",
            title: "会員登録",
            url: "/register",
            log_action: "view_register",
            views: {
                "@": {
                    templateUrl: "app/register/register.html",
                    controller: "RegisterController",
                    controllerAs: "reg"
                }
            },
            resolve: {
                init: ["User", function(User) {
                    return User.initModel()
                }
                ]
            }
        }).state("register.complete", {
            title: "会員登録完了",
            url: "/complete",
            log_action: "view_register_complete",
            views: {
                "@": {
                    templateUrl: "app/register/complete.html",
                    controller: "RegisterController",
                    controllerAs: "reg"
                }
            }
        }).state("faq", {
            parent: "app",
            title: "よくある質問",
            url: "/faq",
            log_action: "view_faq",
            views: {
                "@": {
                    templateUrl: "app/static/faq.html"
                }
            }
        }).state("company", {
            parent: "app",
            title: "会社概要",
            url: "/company",
            log_action: "view_company",
            views: {
                "@": {
                    templateUrl: "app/static/company.html"
                }
            }
        }).state("contact", {
            parent: "app",
            title: "お問い合わせ",
            url: "/contact",
            log_action: "view_contact",
            views: {
                "@": {
                    templateUrl: "app/static/contact.html"
                }
            }
        }).state("policy", {
            parent: "app",
            title: "プライバシーポリシー",
            url: "/privacypolicy",
            log_action: "view_privacypolicy",
            views: {
                "@": {
                    templateUrl: "app/static/policy.html"
                }
            }
        }).state("terms", {
            parent: "app",
            title: "利用規約",
            url: "/terms",
            log_action: "view_terms",
            views: {
                "@": {
                    templateUrl: "app/static/terms.html"
                }
            }
        }).state("member_terms", {
            parent: "app",
            title: "会員規約",
            url: "/member_terms",
            log_action: "view_member_terms",
            views: {
                "@": {
                    templateUrl: "app/static/member_terms.html"
                }
            }
        }).state("links", {
            parent: "app",
            title: "リンク",
            url: "/links",
            log_action: "view_links",
            views: {
                "@": {
                    templateUrl: "app/static/links.html"
                }
            }
        }).state("error", {
            parent: "app",
            url: "/error",
            log_action: "view_error",
            templateUrl: "app/error/500.html"
        }).state("maintenance", {
            parent: "app",
            templateUrl: "app/error/maintenance.html"
        }).state("notfound", {
            parent: "app",
            url: "/:path",
            log_action: "view_404",
            templateUrl: "app/error/404.html"
        }),
        $urlRouterProvider.otherwise("/")
    }
    ])
}
.call(this),
function() {
    angular.module("hibikiWeb").constant("moment", moment).constant("isIE", function() {
        var ieVersion, isIE, ua;
        return ua = window.navigator.userAgent.toLowerCase(),
        isIE = !1,
        ieVersion = null,
        ua.match(/(msie|MSIE)/) || ua.match(/(t|T)rident/) ? (isIE = !0,
        ieVersion = ua.match(/((msie|MSIE)\s|rv:)([\d\.]+)/)[3],
        ieVersion = parseInt(ieVersion)) : isIE = !1,
        {
            isIE: isIE,
            ieVersion: ieVersion
        }
    }).constant("hasFlash", function() {
        var e, fo, hasFlash;
        hasFlash = !1;
        try {
            fo = new ActiveXObject("ShockwaveFlash.ShockwaveFlash"),
            fo && (hasFlash = !0)
        } catch (error) {
            e = error,
            navigator.mimeTypes && navigator.mimeTypes["application/x-shockwave-flash"] && navigator.mimeTypes["application/x-shockwave-flash"].enabledPlugin && (hasFlash = !0)
        }
        return hasFlash
    }).constant("isSafari", function() {
        var browserText, e, isSafari, userAgent;
        isSafari = !1;
        try {
            userAgent = window.navigator.userAgent.toLowerCase(),
            browserText = "",
            browserText = -1 !== userAgent.indexOf("msie") || -1 !== userAgent.indexOf("trident") ? "ie" : -1 !== userAgent.indexOf("edge") ? "edge" : -1 !== userAgent.indexOf("chrome") ? "chrome" : -1 !== userAgent.indexOf("safari") ? "safari" : -1 !== userAgent.indexOf("firefox") ? "firefox" : -1 !== userAgent.indexOf("opera") ? "opera" : "",
            console.log(window.navigator.userAgent),
            console.log(browserText),
            "safari" === browserText && (isSafari = !0)
        } catch (error) {
            e = error,
            isSafari = !1
        }
        return isSafari
    }).constant("platform", function() {
        var bot_ua, isAndroid, isBot, isiOS, isiPad, isiPhone, isiPod, os, osv, platform, ua, v;
        return ua = window.navigator.userAgent.toLowerCase(),
        bot_ua = ["applebot", "baiduspider", "bingbot", "bingpreview", "developers.google.com", "embedly", "googlebot", "gigabot", "hatena::useragent", "ia_archiver", "linkedinbot", "madridbot", "msnbot", "rogerbot", "outbrain", "slackbot", "showyoubot", "yahoo! slurp", "Y!J-", "yandex", "yeti", "yodaobot", "facebookexternalhit", "twitterbot"],
        isBot = RegExp("" + bot_ua.join("|")).test(ua),
        isiPhone = /iphone/.test(ua),
        isiPod = /ipod/.test(ua),
        isiPad = /ipad/.test(ua),
        isiOS = isiPhone || isiPod || isiPad,
        isAndroid = /android/i.test(ua),
        isiOS || isiPad ? (os = "ios",
        v = navigator.appVersion.match(/OS (\d+)_(\d+)_?(\d+)?/),
        osv = parseFloat([v[1], v[2], v[3] || "0"].join("."))) : isAndroid ? (os = "android",
        osv = parseFloat(ua.slice(ua.indexOf("android") + 8))) : (os = "other",
        osv = ""),
        platform = {},
        platform.tablet = isiPad || isAndroid && !/mobile/.test(ua),
        platform.mobile = (isiOS || isAndroid) && !platform.tablet,
        platform.pc = !platform.mobile && !platform.tablet,
        platform.os = os,
        platform.osv = osv,
        platform.bot = isBot,
        platform
    }).constant("videoSettings", {
        options: {
            pc: {
                html5: {
                    hls: {
                        withCredentials: !0
                    }
                },
                hls: {
                    withCredentials: !0
                },
                allowsDvr: !0,
                flash: {
                    swf: "/assets/logica-player/logica-player.swf"
                },
                autoplay: !1,
                preload: "none"
            },
            mobile: {
                html5: {
                    hls: {
                        withCredentials: !0
                    }
                },
                hls: {
                    withCredentials: !0
                },
                allowsDvr: !0,
                autoplay: !1,
                preload: "none"
            }
        },
        template: '<video class="logica-player" id="{{player_id}}" width="100%" height="100%" style="width: 100%; height: 100%;" muted title="{{title}}" preload="none">\n  <p class="pgp-no-js">動画を再生するためにはJavaScriptを有効にしてください。</p>\n</video>',
        programPlayerTemplate: '<video class="logica-player" id="program-player" width="100%" height="100%" style="width: 100%; height: 100%;" muted  preload="none" title="{{title}}" ></video>'
    }).constant("prefectures", [{
        value: 1,
        label: "北海道",
        group: "北海道"
    }, {
        value: 2,
        label: "青森県",
        group: "東北"
    }, {
        value: 3,
        label: "岩手県",
        group: "東北"
    }, {
        value: 4,
        label: "宮城県",
        group: "東北"
    }, {
        value: 5,
        label: "秋田県",
        group: "東北"
    }, {
        value: 6,
        label: "山形県",
        group: "東北"
    }, {
        value: 7,
        label: "福島県",
        group: "東北"
    }, {
        value: 8,
        label: "茨城県",
        group: "関東"
    }, {
        value: 9,
        label: "栃木県",
        group: "関東"
    }, {
        value: 10,
        label: "群馬県",
        group: "関東"
    }, {
        value: 11,
        label: "埼玉県",
        group: "関東"
    }, {
        value: 12,
        label: "千葉県",
        group: "関東"
    }, {
        value: 13,
        label: "東京都",
        group: "関東"
    }, {
        value: 14,
        label: "神奈川県",
        group: "関東"
    }, {
        value: 15,
        label: "新潟県",
        group: "中部"
    }, {
        value: 16,
        label: "富山県",
        group: "中部"
    }, {
        value: 17,
        label: "石川県",
        group: "中部"
    }, {
        value: 18,
        label: "福井県",
        group: "中部"
    }, {
        value: 19,
        label: "山梨県",
        group: "中部"
    }, {
        value: 20,
        label: "長野県",
        group: "中部"
    }, {
        value: 21,
        label: "岐阜県",
        group: "中部"
    }, {
        value: 22,
        label: "静岡県",
        group: "中部"
    }, {
        value: 23,
        label: "愛知県",
        group: "中部"
    }, {
        value: 24,
        label: "三重県",
        group: "近畿"
    }, {
        value: 25,
        label: "滋賀県",
        group: "近畿"
    }, {
        value: 26,
        label: "京都府",
        group: "近畿"
    }, {
        value: 27,
        label: "大阪府",
        group: "近畿"
    }, {
        value: 28,
        label: "兵庫県",
        group: "近畿"
    }, {
        value: 29,
        label: "奈良県",
        group: "近畿"
    }, {
        value: 30,
        label: "和歌山県",
        group: "近畿"
    }, {
        value: 31,
        label: "鳥取県",
        group: "中国"
    }, {
        value: 32,
        label: "島根県",
        group: "中国"
    }, {
        value: 33,
        label: "岡山県",
        group: "中国"
    }, {
        value: 34,
        label: "広島県",
        group: "中国"
    }, {
        value: 35,
        label: "山口県",
        group: "中国"
    }, {
        value: 36,
        label: "徳島県",
        group: "四国"
    }, {
        value: 37,
        label: "香川県",
        group: "四国"
    }, {
        value: 38,
        label: "愛媛県",
        group: "四国"
    }, {
        value: 39,
        label: "高知県",
        group: "四国"
    }, {
        value: 40,
        label: "福岡県",
        group: "九州"
    }, {
        value: 41,
        label: "佐賀県",
        group: "九州"
    }, {
        value: 42,
        label: "長崎県",
        group: "九州"
    }, {
        value: 43,
        label: "熊本県",
        group: "九州"
    }, {
        value: 44,
        label: "大分県",
        group: "九州"
    }, {
        value: 45,
        label: "宮崎県",
        group: "九州"
    }, {
        value: 46,
        label: "鹿児島県",
        group: "九州"
    }, {
        value: 47,
        label: "沖縄県",
        group: "沖縄"
    }]).constant("gender", [{
        value: 0,
        label: "未選択"
    }, {
        value: 1,
        label: "男"
    }, {
        value: 2,
        label: "女"
    }]).constant("allow_mail", [{
        value: !0,
        label: "受け取る"
    }, {
        value: !1,
        label: "受け取らない"
    }]).constant("store_url", {
        ios: "https://itunes.apple.com/jp/app/id843980309",
        android: "https://play.google.com/store/apps/details?id=com.bushiroad.HibikiRadio"
    }).constant("cookieKeys", {
        sessionid: "Hibiki-sessionid",
        uuid: "Hibiki-uuid"
    })
}
.call(this),
function() {
    angular.module("hibikiWeb").config(["$provide", "$locationProvider", "$logProvider", "isIE", "moment", "platform", "cfpLoadingBarProvider", "$httpProvider", "CacheFactoryProvider", "growlNotificationsProvider", "$sceDelegateProvider", "$compileProvider", "$analyticsProvider", function($provide, $locationProvider, $logProvider, isIE, moment, platform, cfpLoadingBarProvider, $httpProvider, CacheFactoryProvider, growlNotificationsProvider, $sceDelegateProvider, $compileProvider, $analyticsProvider) {
        var _ie;
        return $provide.decorator("$uiViewScroll", function() {
            return function(uiViewElement) {
                return angular.element("body.pc #layout header").removeClass("fix").css({
                    bottom: ""
                }),
                angular.element("body.pc #layout .view-container").removeClass("fix").css({
                    width: "",
                    bottom: ""
                }),
                window.scrollTo(0, 0)
            }
        }),
        _ie = isIE(),
        "undefined" == typeof window.history || _ie.isIE && _ie.ieVersion <= 9 ? ($locationProvider.html5Mode(!1),
        $locationProvider.hashPrefix("!")) : $locationProvider.html5Mode(!0),
        $logProvider.debugEnabled(!0),
        cfpLoadingBarProvider.includeSpinner = !1,
        moment.locale("ja", {
            weekdays: ["日曜日", "月曜日", "火曜日", "水曜日", "木曜日", "金曜日", "土曜日"],
            weekdaysShort: ["日", "月", "火", "水", "木", "金", "土"]
        }),
        $httpProvider.defaults.withCredentials = !0,
        $httpProvider.defaults.headers.post["Content-Type"] = "application/x-www-form-urlencoded;application/json;charset=utf-8",
        $httpProvider.interceptors.push(["$q", "$rootScope", "$location", "$injector", function($q, $rootScope, $location, $injector) {
            return {
                request: function(config) {
                    var apiHost, user;
                    return apiHost = $injector.get("apiHost"),
                    config.url.match(new RegExp(apiHost)) && (config.headers["X-Requested-With"] = "XMLHttpRequest",
                    angular.isDefined($rootScope.User) && $rootScope.User.loggedIn() && (user = $rootScope.User.get(),
                    config.headers["X-Hibiki-User-Id"] = user.id,
                    config.headers["X-Hibiki-Access-Token"] = user.access_token)),
                    config || $q.when(config)
                },
                response: function(response) {
                    var $http, $timeout, access_token_expires;
                    return access_token_expires = response.headers("X-Hibiki-Access-Token-Expires"),
                    !response.config.cached && angular.isDefined($rootScope.User) && $rootScope.User.loggedIn() && access_token_expires && $rootScope.User.set({
                        expired_at: access_token_expires
                    }),
                    $http = $injector.get("$http"),
                    $timeout = $injector.get("$timeout"),
                    $rootScope = $injector.get("$rootScope"),
                    $http.pendingRequests.length < 1 && $timeout(function() {
                        return $http.pendingRequests.length < 1 ? $rootScope.htmlReady() : void 0
                    }, 700),
                    response || $q.when(response)
                },
                responseError: function(rejection) {
                    var status;
                    return status = rejection.status,
                    $q.reject(rejection)
                }
            }
        }
        ]),
        growlNotificationsProvider.setOptions({
            ttl: 1e4
        }),
        $compileProvider.aHrefSanitizationWhitelist(/^\s*(https?|ftp|mailto|tel|file|intent|jphibikiradio):/),
        $sceDelegateProvider.resourceUrlWhitelist(["self", "https://hibikiradiovms.blob.core.windows.net/static/**", "https://s3-ap-northeast-1.amazonaws.com/hibiki-terms/**", "jphibikiradio://**"])
    }
    ])
}
.call(this),
angular.module("hibikiWeb").run(["$templateCache", function($templateCache) {
    $templateCache.put("app/error/404.html", '<div id="notfound"><div class="page-title">ご指定のページがみつかりません</div><div class="inner"><div class="box"><p style="line-height: 18px;">ご指定のページは一時的にアクセスができない状況にあるか、<br>移動もしくは削除された可能性があります。</p><br><a ui-sref="home">トップページにもどる</a></div></div></div>'),
    $templateCache.put("app/error/500.html", '<div id="error"><div class="page-title">エラーが発生しました</div><div class="inner"><div class="box"><p style="line-height: 18px;">エラーが発生しました。<br>しばらくしてから再度アクセスをお試しください。</p><br><a ui-sref="home">トップページにもどる</a></div></div></div>'),
    $templateCache.put("app/error/maintenance.html", '<div id="maintenance"><header class="page-title">メンテナンス</header><div class="inner"><div class="box"><div><p ng-bind-html="$root.maintenance_reason"></p></div><div class="mt10" ng-if="$root.maintenance_url"><a ng-href="{{$root.maintenance_url}}" class="btn more">詳しくはこちらをご確認ください</a></div></div></div></div>'),
    $templateCache.put("app/favorites/favorites.html", '<div id="favorites"><header class="page-title">お気に入り一覧</header><div class="inner"><div class="menu" ng-show="fav.programs.length > 0"><div class="edit"><a href="#" ng-click="fav.editable = !fav.editable" ng-class="{editing:fav.editable}"><i class="fa fa-check-circle-o fa-fw fa-lg"></i>お気に入りを削除</a></div><div class="autoplay"><i class="fa fa-play-circle-o fa-fw fa-lg"></i>このページで再生を行うとお気に入りに登録された番組を連続で再生します</div></div><div class="favorites" infinite-scroll="fav.loadMore()" infinite-scroll-disabled="fav.busy" infnite-scroll-distance="1"><div class="nocontent" ng-if="fav.programs.length == 0"><p>お気に入り番組が登録されていません。<br>番組詳細からお気に入り番組を登録することができます。</p></div><div class="program" ng-repeat="program in fav.programs"><div class="delete-favorite" ng-show="fav.editable"><a href="#" ng-click="fav.deleteFavorite($event, program.id)" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'user_favorite_unregister\',label:program.id}"><i class="fa fa-minus-circle fa-lg"></i></a></div><div class="inner-right"><div class="table"><div class="thumbnail"><div class="img" ng-style="{\'background-image\': \'url({{program.sp_image_url}})\'}"><a ui-sref="program.detail({access_id:program.access_id})"></a></div></div><div class="info"><div class="title">{{program.name}}</div><div class="latest-episode-name">{{program.latest_episode_name}}</div><div class="cast">{{program.cast}}</div><div class="episode-updated-at"><span class="new-program" ng-show="program.new_program_flg">新番組</span> <span class="new" ng-show="program.update_flg">NEW</span> <span class="updated-date">{{program.episode_updated_at | moment: \'YYYY.MM.DD\'}}</span></div></div></div></div><div class="inner-left"><div class="control"><div class="play"><a href="#" ng-click="fav.playWithQueue($event, program.id)" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'play_user_favorite\',label:program.episode.id}"><div class="table"><div class="cell icon"><span class="play-icon"></span></div><div class="cell"><span>連続再生</span></div></div></a></div><div class="view"><a ui-sref="program.detail({access_id:program.access_id})" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'user_favorite_program\',label:program.id}"><div class="table"><div class="cell icon"><i class="program-icon fa fa-file-o fa-fw"></i></div><div class="cell">番組ページ</div></div></a></div></div></div></div></div></div></div>'),
    $templateCache.put("app/favorites/favorites.sp.html", '<div id="favorites"><header class="page-title">お気に入り一覧</header><div class="inner"><div class="menu"><div class="edit"><a href="#" ng-click="fav.editable = !fav.editable" ng-class="{editing:fav.editable}"><i class="fa fa-check-circle-o fa-fw fa-lg"></i>お気に入りを編集</a></div></div><div class="favorites" infinite-scroll="fav.loadMore()" infinite-scroll-disabled="fav.busy" infnite-scroll-distance="1"><div class="nocontent" ng-if="fav.programs.length == 0"><p>お気に入り番組が登録されていません。<br>番組詳細からお気に入り番組を登録することができます。</p></div><div class="program" ng-repeat="program in fav.programs"><div class="delete-favorite" ng-show="fav.editable"><a href="#" ng-click="fav.deleteFavorite($event, program.id)" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'user_favorite_unregister\',label:program.id}"><i class="fa fa-minus-circle fa-2x"></i></a></div><div class="thumbnail" ng-style="{\'background-image\': \'url({{::program.sp_image_url}})\'}"></div><div class="info"><div class="title">{{::program.name}}</div><div class="latest-episode-name">{{::program.latest_episode_name}}</div><div class="cast">{{program.cast}}</div><div class="episode-updated-at"><span class="new-program" ng-show="program.new_program_flg">新番組</span> <span class="new" ng-show="program.update_flg">NEW</span><div class="updated-date">{{::program.episode_updated_at}}</div></div></div><div class="control"><div class="play"><a href="#" ng-click="fav.play($event, program.access_id)" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'play_user_favorite\',label:program.episode.id}"><div class="table"><div class="cell icon"><span class="play-icon"></span></div><div class="cell"><span>PLAY</span></div></div></a></div><div class="view"><a ui-sref="program.detail({access_id:program.access_id})" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'user_favorite_program\',label:program.id}"><div class="table"><div class="cell icon"><i class="program-icon fa fa-file-o fa-fw"></i></div><div class="cell">番組ページ</div></div></a></div></div></div></div></div></div>'),
    $templateCache.put("app/footer/footer.html", '<div id="app"><div class="app-icon"><div class="icon-container"><img src="/assets/images/app-icon.png" alt=""></div><div class="text"><p>響 - HiBiKi Radio Station - 専用アプリで<br>もっと響を楽しもう！</p></div></div><div class="qrcode"><div class="ios"><div class="qrimage"><img src="assets/images/iphone-app-qr.gif" alt="App Store"></div><p><i class="fa fa-apple fa-fw"></i>iOS</p></div><div class="android"><div class="qrimage"><img src="assets/images/android-app-qr.gif" alt="Google Play"></div><p><i class="fa fa-android fa-fw"></i>Android</p></div></div></div><div class="footer-inner"><div class="copyrights" ng-bind-html="footer.copyright"></div></div>'),
    $templateCache.put("app/footer/footer.sp.html", '<div id="app" ng-cloak=""><div class="app-icon"><div class="icon-container"><img src="/assets/images/app-icon.png" alt=""></div></div><div class="text"><p>響 - HiBiKi Radio Station - 専用アプリで<br>もっと響を楽しもう！</p></div><div class="dl-link"><div class="ios"><a ng-href="{{footer.store_url.ios}}" target="_blank" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'dl_app_store\'}"><img src="/assets/images/appstore.png" alt="App Storeからダウンロード"></a></div><div class="android"><a ng-href="{{footer.store_url.android}}" target="_blank" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'dl_google_play\'}"><img src="/assets/images/googleplay.png" alt="Google playで手に入れよう"></a></div></div></div><div id="footer" ng-cloak=""><div class="sns"><div class="item cell"><a class="facebook" href="#" socialshare="" socialshare-provider="facebook" socialshare-text="響 - HiBiKi Radio Station -" socialshare-url="https://hibiki-radio.jp" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_share_facebook\'}"><div class="item-inner"><span class="icon"></span> Facebook</div></a></div><div class="item cell"><a class="twitter" ng-href="https://twitter.com/intent/follow?screen_name=HiBiKi_radio" target="_blank" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_follow_twitter\'}"><div class="item-inner"><span class="icon"></span> Twitter<script>!function(d,s,id){var js,fjs=d.getElementsByTagName(s)[0],p=/^http:/.test(d.location)?\'http\':\'https\';if(!d.getElementById(id)){js=d.createElement(s);js.id=id;js.src=p+\'://platform.twitter.com/widgets.js\';fjs.parentNode.insertBefore(js,fjs);}}(document, \'script\', \'twitter-wjs\');</script></div></a></div></div><div class="other-menu-list"><div class="row"><div class="menu-item cell"><a ui-sref="faq" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_faq\'}"><span class="menu-text">よくある質問</span></a></div><div class="menu-item cell"><a ui-sref="contact" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_contact\'}"><span class="menu-text">お問い合わせ</span></a></div></div><div class="row"><div class="menu-item cell"><a ui-sref="policy" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_privacypolicy\'}"><span class="menu-text">プライバシーポリシー</span></a></div><div class="menu-item cell"><a ui-sref="links" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_links\'}"><span class="menu-text">当サイトへのリンクについて</span></a></div></div><div class="row"><div class="menu-item cell"><a ui-sref="terms" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_terms\'}"><span class="menu-text">利用規約</span></a></div><div class="menu-item cell"><a ui-sref="member_terms" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_member_terms\'}"><span class="menu-text">会員規約</span></a></div></div><div class="row"><div class="menu-item cell"><a ui-sref="company" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_company\'}"><span class="menu-text">会社概要</span></a></div><div class="menu-item cell"></div></div></div><div class="banners"><div class="row"><div class="banner cell"><a href="https://bushiroad.com/" target="_blank" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'banner_bushiroad\'}"><img src="/assets/images/bushiroad_banner.jpg" alt="ブシロード"></a></div><div class="banner cell"><a href="https://www.bushiroad-music.com/" target="_blank" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'banner_bushiroad_music\'}"><img src="/assets/images/bushiroad_music.jpg" alt="ブシロードMUSIC"></a></div></div><div class="row"><div class="banner cell"><a href="https://www.bushimo.jp/" target="_blank" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'banner_bushiroad_media\'}"><img src="/assets/images/bushimo.jpg" alt="ブシモ"></a></div><div class="banner cell"><a href="https://bushiroad-ecshop.com/" target="_blank" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'hibiki_ec_shop\'}"><img src="/assets/images/bushiroad_ec_shop_banner.jpeg" alt="ブシロード EC SHOP"></a></div></div></div><div class="jasrac"><div class="jasrac-inner" style="margin:auto"><div class="mark"><img src="/assets/images/jasrac.jpg" alt="JASRAC"></div><div class="jasrac-number"><div>JASRAC許諾番号</div><div>9017659001Y31016</div></div><div class="mark"><img src="/assets/images/elicense_logo.png" alt="JASRAC"></div><div class="jasrac-number"><div>ID36922</div></div></div></div><div class="bushiroad-copyright"><span class="copyright"><i class="fa fa-copyright"></i>bushiroad All Rights Reserved.</span></div></div><div class="copyrights" ng-bind-html="footer.copyright"></div>'),
    $templateCache.put("app/header/header.html", '<div id="menu" side-navbar="" notifications="header.notifications" dailyvoice="header.daily_voice"></div>'),
    $templateCache.put("app/header/header.sp.html", '<div side-navbar="" dailyvoice="header.daily_voice"></div>'),
    $templateCache.put("app/informations/information_detail.html", '<article id="information-detail"><header class="page-title">{{info_d.data.kind_name}}</header><div class="inner"><section class="information-detail"><span class="day">{{info_d.data.day}}</span><h2 class="title">{{info_d.data.name}}</h2><div ng-repeat="p in info_d.data.information_parts" ng-if="!info_d.data.html_description"><div class="info-image" ng-if="p.pc_image_url && !info_d.isMobile"><img ng-src="{{p.pc_image_url}}"></div><div class="info-image" ng-if="p.sp_image_url && info_d.isMobile"><img ng-src="{{p.sp_image_url}}"></div><p ng-bind-html="p.description | linky:\'_blank\'"></p></div><div ng-bind-html="info_d.data.html_description" ng-if="!!info_d.data.html_description"></div><div class="link" ng-if="info_d.data.link_url"><a ng-href="{{info_d.data.link_url}}" target="_blank">詳しくはこちら</a></div></section></div></article>'),
    $templateCache.put("app/informations/informations.html", '<div id="informations"><div class="page-title">お知らせ一覧</div><div class="inner"><div class="information-kind form-group"><label>カテゴリ</label> <label for="information-kind-select" class="select-label"><select id="information-kind-select" name="" class="form-control" ng-model="info.selectedKind" ng-change="info.searchKind()" ng-options="kind.value as kind.label for kind in info.kinds"></select></label></div><div class="infomations-list"><div class="nocontent" ng-show="info.informations.length == 0"><p>お知らせはありません</p></div><ul infinite-scroll="info.loadMore()" infinite-scroll-disabled="info.busy" infinite-scroll-distance="2"><li ng-repeat="i in info.informations"><div class="info"><a ui-sref=".detail({info_id: i.id})" class="info-table"><div class="info-title-cell"><img ng-src="{{i.sp_image_url}}" class="info-title-img" ng-if="i.sp_image_url"><div><div><span class="infolabel {{info.labelClass(i.kind)}}">{{i.kind_name}}</span><span class="publish-start-at">{{i.day}}</span></div><div class="info-title">{{i.name}}</div></div></div></a></div></li></ul></div></div></div>'),
    $templateCache.put("app/instantvoices/instantvoices.html", '<div id="instant-voices"><div class="page-title">インスタントボイス<span class="sub-title">声優陣のヒトコトをボイスでお届け！</span></div><div class="inner"><div class="instant-voice-list" ng-cloak=""><ul infinite-scroll="iv.loadMore()" infinite-scroll-disabled="iv.busy" infinite-scroll-distance="2"><li ng-repeat="voice in iv.instant_voices"><instant-voice voice="voice" star="true" charlength="50" logaction="\'play_instant_voice\'"></instant-voice></li></ul></div></div></div>'),
    $templateCache.put("app/login/login.html", '<div id="login"><header class="page-title">ログイン</header><div class="inner"><div class="login-form"><form id="login-form" name="login_form" ng-submit="login.submit()" novalidate=""><div class="form-group"><label>メールアドレス</label> <input id="" type="email" name="email" class="form-control" autocorrect="off" ng-model="login.email" placeholder="メールアドレス" ng-required="true"><div ng-messages="login_form.email.$error" ng-if="login_form.email.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p><p class="error-message" ng-message="email">メールアドレスを入力してください</p></div></div><div class="form-group"><label>パスワード</label> <input id="" type="password" name="password" class="form-control" autocorrect="off" ng-model="login.password" placeholder="パスワード" ng-required="true" ng-minlength="8" ng-pattern="/^[0-9a-zA-Z-_\\!\\.]+$/"><div ng-messages="login_form.password.$error" ng-if="login_form.password.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p><p class="error-message" ng-message="pattern">半角英数字と記号（-_!.）で入力してください</p><p class="error-message" ng-message="minlength">8文字以上を入力してください</p></div></div><div class="btn-group"><button type="submit" class="btn" ng-disabled="login_form.$invalid">ログイン</button></div><div class="link-group"><ul><li><a ui-sref="register">新規会員登録</a></li><li><a ui-sref="password_reminder">パスワードを忘れた方はこちら</a></li></ul></div></form></div></div></div>'),
    $templateCache.put("app/logout/logout.html", '<div id="logout"><div class="page-title">ログアウト</div><div class="inner"><div class="box"><div><p>ログアウトしました。</p><br><p><span>{{logout.countdown}}秒後にTOPページに移動します</span></p></div></div></div></div>'),
    $templateCache.put("app/main/main.html", '<banner-carousel banners="main.banners"></banner-carousel><div class="container"><program-list></program-list><div class="adsense-middle"><div class="ad-cont"><div class="adsense-block" ng-include="\'app/components/adsense/responsive-rectangle.html\'"></div></div><div class="ad-cont"><div class="adsense-block" ng-include="\'app/components/adsense/responsive-rectangle.html\'"></div></div></div><div class="other-content"><div class="twitter-widget" ng-cloak=""><div class="title"><span class="twitter-icon"><i class="fa fa-twitter"></i></span>@HiBiKi_radioさんのツイート</div><div><a class="twitter-timeline" href="https://twitter.com/HiBiKi_radio" data-widget-id="649140381309755392" data-chrome="noheader">@HiBiKi_radioさんのツイート</a></div></div><div class="twitter-widget" ng-cloak=""><div class="title"><span class="twitter-icon"><i class="fa fa-twitter"></i></span>#響ラジオのツイート</div><div><a class="twitter-timeline" href="https://twitter.com/hashtag/%E9%9F%BF%E3%83%A9%E3%82%B8%E3%82%AA" data-widget-id="630973343001264128" data-chrome="noheader" lang="ja">#響ラジオ のツイート</a></div></div></div></div>'),
    $templateCache.put("app/main/main.sp.html", '<sp-banner-carousel banners="main.banners"></sp-banner-carousel><div class="news-top5" ng-if="main.news.length > 0" ng-cloak=""><div class="title">最新情報</div><slick dots="true" autoplay="false" speed="600" arrows="false"><div class="news-container" ng-repeat="n in main.news"><news data="n"></news></div></slick><div><a class="more" ui-sref="notifications">もっと見る</a></div></div><div class="container"><program-list programs="main.programs"></program-list><div class="adsense-middle"><div class="ad-cont"><div class="adsense-block" ng-include="\'app/components/adsense/responsive-rectangle.html\'"></div></div></div><div class="twitter-widget" ng-cloak=""><div class="title">Twitter<span class="twitter-icon"><i class="fa fa-twitter"></i></span></div><div><a class="twitter-timeline" href="https://twitter.com/hashtag/%E9%9F%BF%E3%83%A9%E3%82%B8%E3%82%AA" data-widget-id="630973343001264128" data-chrome="noheader" lang="ja">#響ラジオ のツイート</a></div></div><div class="adsense-footer"><div class="ad-cont"><div class="adsense-block" ng-include="\'app/components/adsense/responsive-rectangle.html\'"></div></div></div></div>'),
    $templateCache.put("app/mypage/mypage.html", '<div id="mypage"><div class="page-title">マイページ</div><div class="inner"><div class="box"><div class="table"><div class="submenu"><ul><li><a ui-sref="mypage" ng-class="{active: $root.currentState == \'mypage\'}" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'mypage_profile\'}"><span>会員情報</span></a></li><li><a ui-sref=".profile_edit" ng-class="{active: $root.currentState == \'mypage.profile_edit\'}" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'mypage_profile_edit\'}"><span>会員情報変更</span></a></li><li><a ui-sref=".password" ng-class="{active: $root.currentState == \'mypage.password\'}" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'mypage_password_change\'}"><span>パスワード変更</span></a></li><li><a ui-sref=".withdraw" ng-class="{active: $root.currentState == \'mypage.withdraw\'}" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'mypage_withdraw\'}"><span>退会</span></a></li></ul></div><div class="user-info"><div ui-view=""><table class="table"><tr><td>メールアドレス</td><td>{{my.regist_user.email}}</td></tr><tr><td>ニックネーム</td><td>{{my.regist_user.name}}</td></tr><tr><td>都道府県</td><td>{{my.regist_user.prefecture_id|prefectures_label}}</td></tr><tr><td>生年月日</td><td>{{my.regist_user.birth_date|moment:\'YYYY年MM月DD日\'}}</td></tr><tr><td>性別</td><td>{{my.regist_user.gender|gender_label}}</td></tr><tr><td>お知らせメール</td><td>{{my.regist_user.allow_mail_flg|allow_mail_label}}</td></tr></table></div></div></div></div></div></div>'),
    $templateCache.put("app/notifications/notifications.html", '<div id="notifications"><div class="page-title">最新情報</div><div class="inner"><div class="notifications-list"><ul infinite-scroll="nt.loadMore()" infinite-scroll-disabled="nt.busy" infinite-scroll-distance="2"><li ng-repeat="n in nt.notifications"><news data="n" large="true"></news></li></ul></div></div></div>'),
    $templateCache.put("app/notifications/notifications.sp.html", '<div id="notifications"><div class="page-title">最新情報</div><div class="inner"><div class="notifications-list"><ul><li ng-repeat="n in nt.notifications"><news data="n" large="true"></news></li></ul></div></div></div>'),
    $templateCache.put("app/password_reminder/password_reminder.html", '<div id="password-reminder"><header class="page-title">パスワードリマインダー</header><div class="inner"><div class="box"><div class="description"><p>お客様のアカウントを確認します。</p><p>ログインに使用しているメールアドレスを入力して送信ボタンを押してください</p></div><form id="password-reminder-form" name="password_reminder_form" ng-submit="pr.submit()"><div class="form-group"><label>メールアドレス</label> <input class="form-control" type="email" name="email" autocorrect="off" ng-model="pr.email" placeholder="メールアドレス" ng-required="true"><div ng-messages="password_reminder_form.email.$error" ng-if="password_reminder_form.email.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p><p class="error-message" ng-message="email">メールアドレスを入力してください</p></div></div><div class="btn-group"><button type="submit" class="btn" ng-disabled="password_reminder_form.$invalid">送信</button></div><div class="link-group" ng-if="!$root.standalone"><ul><li><a ui-sref="register">新規会員登録</a></li><li><a ui-sref="login">ログイン</a></li></ul></div></form></div></div></div>'),
    $templateCache.put("app/program/program.html", '<div id="program-image"><img class="program-image program-image-animate" ng-src="{{::program.detail.pc_image_url}}" alt="{{::program.detail.name}}" ng-if="program.detail.pc_image_url"><div class="program-image-instead" ng-if="!program.detail.pc_image_url">{{::program.detail.name}}</div><div class="program-play-wrapper"><div class="program-play" ng-click="program.play($event)" ng-if="program.video.id" logicalog-on="click" logicalog-properties="{category:\'tap\',action:program.playing() ? \'pause_episode\' : \'play_episode\', label: program.detail.episode.id + \'-0\'}" ng-cloak=""><div class="program-play-container"><div><span class="play-icon" ng-class="{play:program.playing(), loading:program.loading}"></span></div><div><div class="updated-at">{{::program.detail.episode_updated_at | moment: "YYYY.MM.DD"}} UP</div><div class="program-name">{{::program.detail.episode.name}}を聞く</div></div></div></div><div class="program-backstage-play" ng-click="program.play($event, true)" ng-if="program.additional_video.id" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'play_episode\',label: program.detail.episode.id + \'-1\'}" logicalog-if="$root.User.loggedIn()" ng-cloak=""><div class="program-play-container"><div><span class="play-icon" ng-class="{play:program.playing(true), loading:program.backstage_loading}"></span></div><div><div class="updated-at">{{::program.detail.episode_updated_at | moment: "YYYY.MM.DD"}} UP</div><div class="program-name">楽屋裏を聞く</div></div></div></div></div></div><div id="program-menu"><ul class="menu"><li><a ng-class="{active: program.stateIs(\'detail\')}" ui-sref=".detail({access_id: program.detail.access_id})" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_detail\',label:program.detail.id}">番組詳細</a></li><li><a ng-class="{active: program.stateIs(\'introductions\')}" ui-sref=".introductions({access_id: program.detail.access_id})" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_segments\',label:program.detail.id}">コーナー紹介</a></li><li ng-if="program.products.length > 0"><a ng-class="{active: program.stateIs(\'products\')}" ui-sref=".products({access_id: program.detail.access_id})" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_product_informations\',label:program.detail.id}">商品</a></li><li ng-if="program.informations.length > 0"><a ng-class="{active: program.stateIs(\'informations\')}" ui-sref=".informations({access_id: program.detail.access_id})" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_program_informations\',label:program.detail.id}">お知らせ<span class="info-count" ng-cloak="">{{program.informations.length}}</span></a></li></ul><a href="mailto:{{program.detail.email}}" target="_blank" class="put-letter" ng-if="program.detail.email" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_mail\',label:program.detail.id}"><i class="icon fa fa-envelope-o fa-lg"></i>{{program.detail.email}}</a> <a href="#" ng-click="program.letter(program.detail.message_form_url)" class="put-letter" ng-if="program.detail.message_form_url" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_message\',label:program.detail.id}"><i class="icon fa fa-file-text-o fa-lg"></i>お便りを出す</a></div><div id="put-letter"><a href="#" ng-click="program.letter(program.detail.message_form_url)" class="put-letter" ng-if="program.detail.message_form_url" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_message\',label:program.detail.id}"><i class="icon fa fa-file-text-o fa-lg"></i>お便りを出す</a><div id="mailto" ng-if="program.detail.email"><div class="caption">直接メールを送りたい方はこちら</div><a href="mailto:{{program.detail.email}}" target="_blank" class="put-letter" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_mail\',label:program.detail.id}"><i class="icon fa fa-envelope-o fa-lg"></i>{{program.detail.email}}</a></div></div><div class="container"><div ui-view=""></div></div>'),
    $templateCache.put("app/program/program.sp.html", '<div id="program-image" ng-cloak=""><img class="program-image program-image-animate" ng-src="{{::program.detail.sp_image_url}}" alt="{{::program.detail.name}}" ng-if="program.detail.sp_image_url"><div class="program-image-instead" ng-if="!program.detail.sp_image_url">{{::program.detail.name}}</div><div class="program-play-wrapper"><div class="program-play" ng-click="program.play($event, null, true)" ng-if="program.detail.episode.video.id" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'play_episode\',label:program.detail.id,value:program.detail.episode.id}" ng-cloak=""><div class="program-play-container"><div><span class="play-icon" ng-class="{play:program.playing()}"></span></div><div><div class="updated-at">{{::program.detail.episode.updated_at | moment: "YYYY.MM.DD"}} UP</div><div class="program-name">{{::program.detail.episode.name}}を聞く</div></div></div></div><div class="program-backstage-play" ng-click="program.play($event, true, true)" ng-if="program.detail.episode.additional_video.id" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'play_episode_appendix\',label:program.detail.id,value:program.detail.episode.id}" ng-cloak=""><div class="program-play-container"><div><span class="play-icon" ng-class="{play:program.playing(true)}"></span></div><div><div class="updated-at">{{::program.detail.episode.updated_at | moment: "YYYY.MM.DD"}} UP</div><div class="program-name">楽屋裏を聞く</div></div></div></div></div></div><div id="program-menu" ng-cloak=""><ul class="menu"><li><a ng-class="{active: program.stateIs(\'detail\')}" ui-sref=".detail({access_id: program.detail.access_id})" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_detail\',label:program.detail.id}">番組詳細</a></li><li><a ng-class="{active: program.stateIs(\'introductions\')}" ui-sref=".introductions({access_id: program.detail.access_id})" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_segments\',label:program.detail.id}">コーナー紹介</a></li><li ng-if="program.products.length > 0"><a ng-class="{active: program.stateIs(\'products\')}" ui-sref=".products({access_id: program.detail.access_id})" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_product_informations\',label:program.detail.id}">商品</a></li><li ng-if="program.informations.length > 0"><a ng-class="{active: program.stateIs(\'informations\')}" ui-sref=".informations({access_id: program.detail.access_id})" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_program_informations\',label:program.detail.id}">お知らせ<span class="info-count" ng-cloak="">{{program.informations.length}}</span></a></li></ul></div><div class="container"><div ui-view=""></div></div>'),
    $templateCache.put("app/register/complete.html", '<div id="register-complete"><header class="page-title">メールアドレス確認メールを送信しました。</header><div class="inner"><div class="box"><p>{{reg.user.email}}宛にメールアドレス確認メールを送信しました。</p><p>メールに記載されているURLをクリックして会員登録を完了してください。</p><div class="mt30"><a ui-sref="login">ログイン</a></div></div></div></div>'),
    $templateCache.put("app/register/member_terms.html", '<div class="dialog"><iframe class="outer-html" ng-src="{{$root.member_terms_url}}" frameborder="0"></iframe></div>'),
    $templateCache.put("app/register/register.html", '<div id="register"><header class="page-title">会員登録</header><div class="inner"><div class="register-form"><form id="register-form" name="register_form" ng-submit="reg.submit()"><div class="register-text"><ul><li>会員登録は無料です</li><li>簡単な会員登録で下記の特典機能が楽しめます。</li><li><ul class="list"><li>・お好きな番組を「お気に入り登録」出来ます。</li><li>・「お気に入り登録をした番組」にて、「お気に入り一覧」より連続再生ができます。</li><li>・「チャプター頭出し機能」にて、より便利にラジオをお聴き頂けます。（※1）</li><li>・「お知らせ情報」をメールにて受け取る事が出来ます。</li><li>・その他、特典機能を追加予定です。</li></ul></li></ul><div class="caption">※1 一部の番組が対象となります。</div></div><div class="form-group"><label>メールアドレス</label> <input type="email" name="email" class="form-control" autocorrect="off" ng-model="reg.user.email" placeholder="ログインに使用するメールアドレスを入力してください" ng-required="true"><div ng-messages="register_form.email.$error" ng-if="register_form.email.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p><p class="error-message" ng-message="email">メールアドレスを入力してください</p></div></div><div class="form-group"><label>ニックネーム</label> <input type="text" name="name" class="form-control" ng-model="reg.user.name" placeholder="ログイン後に表示される名前を入力してください" ng-required="true"><div ng-messages="register_form.name.$error" ng-if="register_form.name.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p></div></div><div class="row"><div class="col"><div class="form-group"><label>パスワード</label> <input id="" type="password" name="password" class="form-control" autocorrect="off" ng-model="reg.user.password" placeholder="8文字以上の半角英数字と記号（-_!.）を入力してください" ng-required="true" ng-minlength="8" ng-maxlength="32" ng-pattern="/^[0-9a-zA-Z-_\\!\\.]+$/"><div ng-messages="register_form.password.$error" ng-if="register_form.password.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p><p class="error-message" ng-message="pattern">半角英数字と記号（-_!.）で入力してください</p><p class="error-message" ng-message="minlength">8文字以上で入力してください</p><p class="error-message" ng-message="maxlength">32文字以下で入力してください</p></div></div></div><div class="col"><div class="form-group"><label>パスワード(確認)</label> <input id="" type="password" name="password_confirm" class="form-control" autocorrect="off" ng-model="reg.user.password_confirm" placeholder="確認のため同じパスワードを入力してください" ng-required="true" compare-to="reg.user.password"><div ng-messages="register_form.password_confirm.$error" ng-if="register_form.password_confirm.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p><p class="error-message" ng-message="compareTo">パスワードが一致していません</p></div></div></div></div><div class="row"><div class="col"><div class="form-group"><label>都道府県</label> <label for="prefecture" class="select-label"><select id="prefecture" class="form-control" name="prefecture" ng-model="reg.user.prefecture_id" ng-options="option.value as option.label group by option.group for option in reg.prefectures" ng-required="true"><option value="">都道府県を選択してください</option></select></label><div ng-messages="register_form.prefecture.$error" ng-if="register_form.prefecture.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p></div></div></div><div class="col"><div class="form-group"><label>生年月日</label><div class="birth-date"><label for="birth-date-year" class="select-label"><select id="birth-date-year" class="form-control" name="birth_date[year]" ng-model="reg.user.birth_date_parts.year" ng-options="option for option in reg.birth_date.year" ng-required="true"><option value=""></option></select></label> <span>年</span> <label for="birth-date-month" class="select-label"><select id="birth-date-month" class="form-control" name="birth_date[month]" ng-model="reg.user.birth_date_parts.month" ng-options="option for option in reg.birth_date.month" ng-required="true"><option value=""></option></select></label> <span>月</span> <label for="birth-date-day" class="select-label"><select id="birth-date-day" class="form-control" name="birth_date[day]" ng-model="reg.user.birth_date_parts.day" ng-options="option for option in reg.birth_date.day"><option value=""></option></select></label> <span>日</span></div><input type="hidden" name="birth_date" ng-model="reg.user.birth_date" ng-required="true"><div ng-messages="register_form.birth_date.$error" ng-if="register_form.birth_date.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p></div></div></div></div><div class="row"><div class="col"><div class="form-group"><label>性別</label> <label><input id="" type="radio" name="gender" ng-model="reg.user.gender" value="1"> 男</label> <label><input id="" type="radio" name="gender" ng-model="reg.user.gender" value="2"> 女</label></div></div><div class="col"><div class="form-group allow-mail"><label>お知らせメールを</label> <label><input id="" type="radio" name="allow_mail_flg" ng-model="reg.user.allow_mail_flg" value="true"> 受け取る</label> <label><input id="" type="radio" name="allow_mail_flg" ng-model="reg.user.allow_mail_flg" value="false"> 受け取らない</label><div class="caption">響 -Hibiki Radio Station-からのオススメ情報をお届けします</div><div ng-messages="register_form.allow_mail_flg.$error" ng-if="register_form.allow_mail_flg.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p></div></div></div></div><div class="link-group"><ul><li class="terms"><a href="#" ng-click="reg.showTerms()">利用規約</a>と<a href="#" ng-click="reg.showTerms(\'member\')">会員規約</a>に同意の上、登録をしてください。</li><li>すでに会員の方は<a ui-sref="login">こちら</a>からログインしてください</li></ul></div><div class="btn-group"><button type="submit" class="btn" ng-disabled="register_form.$invalid">登録</button><div class="mt20"><p>登録ボタンを押すと確認メールを入力したメールアドレスに送信します。</p><p>「@hibiki-radio.jp」からのメールを受信出来るように、事前にメール受信の設定をお願い致します。</p></div></div></form></div></div></div>'),
    $templateCache.put("app/register/terms.html", '<div class="dialog"><iframe class="outer-html" ng-src="{{$root.terms_url}}" frameborder="0"></iframe></div>'),
    $templateCache.put("app/static/company.html", '<div id="company"><div class="page-title">会社概要</div><div class="inner"><div class="box"><table class="table"><tr><td>商号</td><td>株式会社響</td></tr><tr><td>所在地</td><td>〒164-0011 東京都中野区中央1-38-1 住友中野坂上ビル6階</td></tr><tr><td>設立</td><td>2016年9月29日</td></tr><tr><td>資本金</td><td>900万円</td></tr><tr><td>取引銀行</td><td>三井住友銀行新宿西口支店</td></tr><tr><td>代表者</td><td>代表取締役社長 岩倉亜貴</td></tr><tr><td>取締役</td><td>小原正司</td></tr><tr><td>取引先</td><td>株式会社ブシロード</td></tr><tr><td>業務内容</td><td><ul><li>インターネットラジオ番組制作及びポータルサイト「響 -HiBiKi Radio Station-」の運営</li><li>タレントマネージメント</li><li>その他関連業務</li></ul></td></tr><tr><td>従業員数</td><td>16人</td></tr></table><h2 class="mt20">アクセス</h2><div class="access"><iframe src="https://www.google.com/maps/embed?pb=!1m14!1m8!1m3!1d3240.1647175736584!2d139.683582!3d35.697564!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x0%3A0x27ad15cbe4a382df!2z5L2P5Y-L5LiN5YuV55SjKOagqikg5L2P5Y-L5Lit6YeO5Z2C5LiK44OT44Or!5e0!3m2!1sja!2sjp!4v1445998409224" width="100%" height="500" frameborder="0" style="border:0" allowfullscreen="" ng-if="!$root.platform.mobile"></iframe><iframe src="https://www.google.com/maps/embed?pb=!1m14!1m8!1m3!1d3240.1647175736584!2d139.683582!3d35.697564!3m2!1i1024!2i768!4f13.1!3m3!1m2!1s0x0%3A0x27ad15cbe4a382df!2z5L2P5Y-L5LiN5YuV55SjKOagqikg5L2P5Y-L5Lit6YeO5Z2C5LiK44OT44Or!5e0!3m2!1sja!2sjp!4v1445998409224" width="100%" height="300" frameborder="0" style="border:0" allowfullscreen="" ng-if="$root.platform.mobile"></iframe></div></div></div></div>'),
    $templateCache.put("app/static/contact.html", '<div id="contact"><div class="page-title">お問い合わせ</div><div class="inner"><div class="box">HibiKi Radio Stationへのご連絡は下記のメールアドレスよりお送りください。<div class="mailto"><a href="mailto:r-info@hibiki-radio.jp">r-info@hibiki-radio.jp</a></div><p><span>【お客様へのお願い】</span><br>お問い合わせの際は必ず、@hibiki-radio.jpからのメールを受信出来る様に、事前にメール受信の設定をお願い致します。</p></div></div></div>'),
    $templateCache.put("app/static/faq.html", '<div id="faq"><div class="page-title">よくある質問</div><div class="inner"><div class="box"><section><div class="question">Q1. 番組は無料で聴けますか？</div><div class="answer"><p>全番組無料でお楽しみいただけます。</p></div></section><section><div class="question">Q2. 会員登録は無料ですか？</div><div class="answer"><p>会員登録は無料です。</p></div></section><section><div class="question">Q3. 会員登録をすると、どんなサービスを受けられますか？</div><div class="answer"><p>当サイト、アプリにて下記の特典機能が楽しめます。</p><ul class="list"><li>お好きな番組を「お気に入り登録」出来ます。</li><li>「お気に入り登録をした番組」にて、「お気に入り一覧」より連続再生ができます。</li><li>「チャプター頭出し機能」にて、より便利にラジオをお聴き頂けます。（※1）</li><li>「お知らせ情報」をメールにて受け取る事が出来ます。</li><li>その他、特典機能を追加予定です。</li></ul><div class="caption">※1 一部の番組が対象となります。</div></div></section><section><div class="question">Q4. 「お気入り登録」はどうやって行うのですか？</div><div class="answer"><p class="color-red">※「お気入り登録機能」を使用するためには、「会員登録」が必要となります。</p><p>「番組詳細」ページに表示されます「お気に入りに追加」というボタンを押して登録を行ってください。</p><p>お気に入り登録された番組は「お気に入り一覧」にて確認・連続再生を行う事が出来ます。</p></div></section><section><div class="question">Q5. 「チャプター機能」はどうやって楽しむのでしょうか？</div><div class="answer"><p class="color-red">※ 「チャプター機能」を使用するためには、「会員登録」が必要となります。</p><p>WEB版の場合:</p><p>「チャプター機能」に対応した番組を再生しますと、再生バーに「▼」マークが表示されます。</p><p>各コーナー毎にチャプター設定されているため、「送りボタン」でお好きなコーナーから再生をする事が出来ます。</p><br><p>アプリの場合:</p><p>「チャプター機能」に対応した番組を再生し、番組画像をタップして頂くと番組画像部分にコーナー情報が表示されます。</p><p>各コーナー毎にチャプター設定されているため、チャプター表示を左右へスライドする事により、お好きなコーナーから再生をする事が出来ます。</p></div></section><section><div class="question">Q6. 会員情報の更新や退会はどこで行えばよいのですか？</div><div class="answer"><p>WEB版の場合:</p><p>会員ログインをして頂いた後、「<a ui-sref="mypage">マイページ</a>」ボタンを押してください。</p><br><p>アプリの場合:</p><p>WEBサイトでの操作になりますので改めて会員ログインをして頂いた後、「<a ui-sref="mypage">マイページ</a>」に進んでください。</p><br><p>「マイページ」にて、会員情報の更新や退会手続きを行うことが出来ます。</p></div></section><section><div class="question">Q7. 番組はどうやってダウンロード出来ますか？</div><div class="answer"><p>原則、全番組ダウンロードは出来ません。</p></div></section><section><div class="question">Q8. 番組が聴けません。</div><div class="answer"><p>・音声・動画の再生につきましては以下をご確認ください。</p><ul><li>Webブラウザの設定でJavascriptが有効になっているか</li><li>Webブラウザの設定でcookieが有効になっているか</li></ul><p>当サイトは以下の環境を推奨しております。</p><br><div><p>【Windows(7以降)をお使いの場合】</p><ul><li>InternetExplorer最新版</li><li>Firefox最新版</li><li>Google Chrome最新版</li></ul><p>※ ブラウザの自動アップデート機能が有効になっていない場合は、最新版へのアップデートを推奨します。</p><p>※ InternetExplorer9,10に関して、動作確認は行っておりますが、ページによっては正しく表示されない場合があります。</p></div><br><div><p>【Mac OSX(10.6以降)をお使いの場合】</p><ul><li>Safari 5.x以降</li><li>Firefox最新版</li><li>Google Chrome最新版</li></ul><p>※ ブラウザの自動アップデート機能が有効になっていない場合は、最新版へのアップデートを推奨します。</p></div><br><div><p>【Android(4.4以降)をお使いの場合】</p><ul><li>Google Chrome</li><li>各OSで標準搭載されているブラウザ</li></ul></div><br><div><p>【iOS(7以降)をお使いの場合】</p><ul><li>各OSで標準搭載されているブラウザ</li></ul></div><br><p>一部のOS、ブラウザ、プレイヤーで正常に再生できないことがございますので、あらかじめご了承ください。</p><p>その場合は他のブラウザにて再生をお試しください。</p><p>iPhone,Androidのブラウザで正常に再生出来ない場合は、アプリでの再生をお試しください。</p></div></section></div></div></div>'),
    $templateCache.put("app/static/links.html", '<div id="link-free"><div class="page-title">当サイトのリンクについて</div><div class="inner"><div class="box"><div><p>当サイトはリンクフリーです。</p><p>リンクは必ずトップページにお願い致します。</p><p>バナーは下記を保存してお使いになるか、指定のリンクをお使いください。</p></div><div class="link-info"><div class="site-name">サイト名: 響 -HiBiKi Radio Station-</div><div class="link-url">リンク先URL: https://hibiki-radio.jp/</div><div class="banner"><a href="assets/images/hibiki-link-logo-s.jpg" target="_blank"><img src="assets/images/hibiki-link-logo-s.jpg" alt=""></a> <span>https://hibiki-radio.jp/assets/images/hibiki-link-logo-s.jpg</span></div><div class="banner"><a href="assets/images/hibiki-link-logo-m.jpg" target="_blank"><img src="assets/images/hibiki-link-logo-m.jpg" alt=""></a> <span>https://hibiki-radio.jp/assets/images/hibiki-link-logo-m.jpg</span></div></div></div></div></div>'),
    $templateCache.put("app/static/member_terms.html", '<div id="member-terms"><header class="page-title">会員規約</header><div class="inner"><div class="box"><iframe class="outer-html" ng-src="{{$root.member_terms_url}}" frameborder="0"></iframe></div></div></div>'),
    $templateCache.put("app/static/policy.html", '<div id="privacy-policy"><div class="page-title">プライバシーポリシー</div><div class="inner"><div class="box"><iframe class="outer-html" ng-src="{{$root.privacy_policy_url}}" frameborder="0"></iframe></div></div></div>'),
    $templateCache.put("app/static/terms.html", '<div id="terms"><header class="page-title">利用規約</header><div class="inner"><div class="box"><iframe class="outer-html" ng-src="{{$root.terms_url}}" frameborder="0"></iframe></div></div></div>'),
    $templateCache.put("app/components/adsense/custom.html", '<ins class="adsbygoogle" style="display:inline-block;width:218px;height:123px" data-ad-client="ca-pub-1085110219277383" data-ad-slot="3696004953"></ins><script>\n(adsbygoogle = window.adsbygoogle || []).push({});\n</script>'),
    $templateCache.put("app/components/adsense/responsive-auto.html", '<ins class="adsbygoogle" style="display:block" data-ad-client="ca-pub-1085110219277383" data-ad-slot="3646279350" data-ad-test="on" data-ad-format="auto"></ins><script>(adsbygoogle = window.adsbygoogle || []).push({});</script>'),
    $templateCache.put("app/components/adsense/responsive-horizontal.html", '<ins class="adsbygoogle" style="display:block" data-ad-client="ca-pub-1085110219277383" data-ad-slot="3646279350" data-ad-test="on" data-ad-format="horizontal"></ins><script>(adsbygoogle = window.adsbygoogle || []).push({});</script>'),
    $templateCache.put("app/components/adsense/responsive-rectangle.html", '<ins class="adsbygoogle" style="display:block" data-ad-client="ca-pub-1085110219277383" data-ad-slot="3646279350" data-ad-test="on" data-ad-format="rectangle"></ins><script>(adsbygoogle = window.adsbygoogle || []).push({});</script>'),
    $templateCache.put("app/components/banner-carousel/banner-carousel.html", '<div id="banner-carousel" ng-show="bc.loaded" ng-mouseover="bc.stop()" ng-mouseleave="bc.start()"><div class="main-banner-area"><div class="main-banner current"><a href="{{bc.main[0].link_url}}" target="_blank" ng-if="bc.main[0].kind == 1 || bc.main[0].kind == 2" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_top_banner\',label:bc.main[0].id}"><img ng-src="{{bc.main[0].pc_image_url}}" alt="{{bc.main[0].name}}"></a> <a ui-sref="informations.detail({info_id: bc.main[0].link_url})" ng-if="bc.main[0].kind == 3" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_top_banner\',label:bc.main[0].id}"><img ng-src="{{bc.main[0].pc_image_url}}" alt="{{bc.main[0].name}}"></a> <a ui-sref="program.detail({access_id: bc.main[0].link_url})" ng-if="bc.main[0].kind == 4" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_top_banner\',label:bc.main[0].id}"><img ng-src="{{bc.main[0].pc_image_url}}" alt="{{bc.main[0].name}}"></a> <a href="#" ng-if="bc.main[0].kind == null"><img ng-src="{{bc.main[0].pc_image_url}}" alt="{{bc.main[0].name}}"></a></div><div class="main-banner next"><a href="{{bc.main[1].link_url}}" target="_blank" ng-if="bc.main[1].kind == 1 || bc.main[1].kind == 2" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_top_banner\',label:bc.main[1].id}"><img ng-src="{{bc.main[1].pc_image_url}}" alt="{{bc.main[1].name}}"></a> <a ui-sref="informations.detail({info_id: bc.main[1].link_url})" ng-if="bc.main[1].kind == 3" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_top_banner\',label:bc.main[1].id}"><img ng-src="{{bc.main[1].pc_image_url}}" alt="{{bc.main[1].name}}"></a> <a ui-sref="program.detail({access_id: bc.main[1].link_url})" ng-if="bc.main[1].kind == 4" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_top_banner\',label:bc.main[1].id}"><img ng-src="{{bc.main[1].pc_image_url}}" alt="{{bc.main[1].name}}"></a> <a href="#" ng-if="bc.main[1].kind == null"><img ng-src="{{bc.main[1].pc_image_url}}" alt="{{bc.main[1].name}}"></a></div></div><div class="thumbnails-area"><div class="thumbnails"><ul class="thumbnail-list"><li ng-repeat="banner in bc.banners"><img ng-src="{{banner.pc_image_url}}" alt="{{banner.name}}"></li></ul></div><div class="thumbnail-frames"><div class="next-frame"></div><div class="prev-btn-area" ng-click="bc.next(\'fast\')"><span class="prev-btn"></span></div><div class="main-frame"></div><div class="prev-frame"></div><div class="next-btn-area" ng-click="bc.prev(\'fast\')"><span class="next-btn"></span></div></div></div></div>'),
    $templateCache.put("app/components/banner-carousel/banner-carousel.sp.html", '<div id="banner-carousel" ng-cloak=""><div class="banner-area"><slick dots="false" autoplay="true" autoplay-speed="5000" speed="600" arrows="false"><div class="banner" ng-repeat="banner in bc.banners"><a href="{{banner.link_url}}" target="_blank" ng-if="banner.kind == 1 || banner.kind == 2" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_top_banner\',label:banner.id}"><img ng-src="{{banner.sp_image_url}}" alt="{{banner.name}}"></a> <a ui-sref="informations.detail({info_id: banner.link_url})" ng-if="banner.kind == 3" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_top_banner\',label:banner.id}"><img ng-src="{{banner.sp_image_url}}" alt="{{banner.name}}"></a> <a ui-sref="program.detail({access_id: banner.link_url})" ng-if="banner.kind == 4" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_top_banner\',label:banner.id}"><img ng-src="{{banner.sp_image_url}}" alt="{{banner.name}}"></a> <a href="#" ng-if="banner.kind == null"><img ng-src="{{banner.sp_image_url}}" alt="{{banner.name}}"></a></div></slick></div></div>'),
    $templateCache.put("app/components/daily-voice/daily-voice.html", '<div class="daily-voice" ng-if="dv.voice"><a href="#" ng-click="dv.playToggle($event)" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'play_daily_voice\',label:dv.voice.id}"><span class="play-icon" ng-hide="dv.playing" ng-class="{loading:dv.loading, play: dv.playing}"></span><div class="now-playing" ng-show="dv.playing"><span class="bar1"></span> <span class="bar2"></span> <span class="bar3"></span> <span class="bar4"></span> <span class="bar5"></span></div><span class="daily-voice-text">本日のつぶやき</span></a></div><div class="daily-voice-title" ng-if="dv.voice"><div ng-class="{playing:dv.playing}"><p>{{dv.voice.name}}</p></div></div>'),
    $templateCache.put("app/components/daily-voice/daily-voice.sp.html", '<div class="daily-voice" ng-class="{play: dv.playing}" ng-if="dv.voice"><a href="#" ng-click="dv.playToggle($event)" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'play_daily_voice\',label:dv.voice.id}"><span class="daily-voice-text">本日の<br>つぶやき</span></a></div>'),
    $templateCache.put("app/components/instant-voice/instant-voice.html", '<div class="instant-voice" ng-cloak=""><div class="thumbnail" ng-style="{\'background-image\': \'url({{iv.voice.actor.pc_image_url || \'/assets/images/no_image.png\'}})\'}"></div><div class="info-wrapper"><div class="info"><span class="actor-name">{{iv.voice.actor.name}}</span> <span class="updated-at">{{iv.voice.updated_at|moment: "MM.DD HH:mm"}}</span></div><div class="description"><div class="description-inner"><p>{{iv.voice.description | truncate: iv.charlength}}</p><div class="play-voice"><a href="#" ng-click="iv.playToggle($event)" logicalog-on="click" logicalog-properties="{category:\'tap\',action:iv.logaction,label:iv.voice.id}"><span class="play-icon" ng-hide="iv.playing" ng-class="{loading:iv.loading, play:iv.playing, member_limited:!$root.User.loggedIn() && iv.voice.member_limited}"></span><div class="now-playing" ng-show="iv.playing"><span class="bar1"></span> <span class="bar2"></span> <span class="bar3"></span> <span class="bar4"></span> <span class="bar5"></span></div></a><div class="duration">{{iv.voice.video.duration|moment_time}}</div></div></div></div></div></div>'),
    $templateCache.put("app/components/instant-voice/main-instant-voice-carousel.html", '<div id="main-instant-voice" ng-show="iv.voices.length > 0"><div class="main-instant-voice-list"><slick dots="false" autoplay="true" autoplay-speed="5000" speed="300" arrows="false" css-ease="linear" pause-on-hover="false" fade="true"><div ng-repeat="voice in iv.voices"><main-instant-voice voice="voice"></main-instant-voice></div></slick></div><div class="more" ng-if="!iv.isMobile"><a ui-sref="instantvoices">もっと見る</a></div></div>'),
    $templateCache.put("app/components/instant-voice/main-instant-voice.html", '<div class="main-instant-voice" ng-show="iv.voice"><div class="actor"><div class="actor-inner"><div class="thumb"><div class="thumb-image" ng-style="{\'background-image\': \'url({{iv.voice.actor.pc_image_url || \'/assets/images/no_image.png\'}})\'}"></div></div><div class="info"><span class="name">{{iv.voice.actor.name | truncate: 8}}</span> <span class="updated-at">{{iv.voice.updated_at|moment: "MM.DD HH:mm"}}</span></div></div></div><div class="play"><a href="#" ng-click="iv.playToggle($event)" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_play_instant_voice_banner\',label:iv.voice.id}"><span class="play-icon" ng-hide="iv.playing" ng-class="{loading:iv.loading, play:iv.playing, member_limited:!$root.User.loggedIn() && iv.voice.member_limited}"></span><div class="now-playing" ng-show="iv.playing"><span class="bar1"></span> <span class="bar2"></span> <span class="bar3"></span> <span class="bar4"></span> <span class="bar5"></span></div></a><div class="duration">{{iv.voice.video.duration|moment_time}}</div></div><div class="message" ng-class="{playing: iv.playing}"><div class="message-inner"><p>{{iv.voice.description}}</p></div></div></div>'),
    $templateCache.put("app/components/instant-voice/main-instant-voice.sp.html", '<div class="main-instant-voice" ng-show="iv.voice"><div class="actor"><div class="thumb" ng-style="{\'background-image\': \'url({{iv.voice.actor.pc_image_url}})\'}"></div></div><div class="message" ng-class="{playing: iv.playing}"><div class="info"><span class="name">{{iv.voice.actor.name}}</span> <span class="updated-at">{{iv.voice.updated_at|moment: "MM.DD HH:mm"}}</span></div><p>{{iv.voice.description}}</p></div><div class="play"><a href="#" ng-click="iv.playToggle($event)" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_play_instant_voice_banner\',label:iv.voice.id}"><span class="play-icon" ng-hide="iv.playing" ng-class="{loading:iv.loading, play:iv.playing, member_limited:!$root.User.loggedIn() && iv.voice.member_limited}"></span><div class="now-playing" ng-show="iv.playing"><span class="bar1"></span> <span class="bar2"></span> <span class="bar3"></span> <span class="bar4"></span> <span class="bar5"></span></div></a><div class="duration">{{iv.voice.video.duration|moment_time}}</div></div></div>'),
    $templateCache.put("app/components/launch-app/launch-app.html", '<div id="launch-app-dialog"><div class="table"><div class="cell"><div class="app-icon"><div class="icon-container"><img src="/assets/images/app-icon.png" alt=""></div></div></div><div class="cell"><p>アプリを利用すると、より響ラジオステーションを楽しむことができます。</p></div></div><div><ul><li><a class="launch-app-link" ng-href="{{appUrl}}"><i class="fa fa-external-link fa-fw"></i> アプリで開く</a></li><li><a ng-href="{{storeUrl}}"><i class="fa fa-external-link fa-fw"></i> アプリをインストールする</a></li><li ng-if="callback"><a href="#" ng-click="closeThisDialog(\'play\')"><i class="fa fa-play fa-fw"></i> このまま再生する</a><p class="caption">※一部の端末では正常に再生されない場合があります。</p></li></ul></div></div>'),
    $templateCache.put("app/components/navbar/navbar.html", '<nav class="menu"><a class="menu-heading" ui-sref="home" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_home\'}"><img src="/assets/images/brand-logo.png" alt="響 Radio Station"></a><div class="navbar-adsense-ad" ng-include="\'app/components/adsense/responsive-rectangle.html\'"></div><daily-voice voice="nb.dailyvoice"></daily-voice><ul class="menu-list"><li class="menu-item"><a ui-sref="home" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_home\'}"><span class="menu-text">トップ</span></a></li><li class="menu-item" ng-show="$root.User.loggedIn()"><a ui-sref="favorites" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_user_favorite\'}"><span class="menu-text">お気に入り一覧</span></a></li><li class="menu-item"><a ui-sref="informations" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_informations\'}"><span class="menu-text">お知らせ一覧</span></a></li></ul><div class="user-menu-list" ng-show="$root.User.loggedIn()"><div class="user"><p class="username"><span class="hello">こんにちは</span><br><span>{{$root.User.get().name}}</span>さん</p></div><div><a class="btn mypage-btn" ui-sref="mypage" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_mypage\'}">マイページ</a></div><div><a class="btn logout-btn" ui-sref="logout" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_logout\'}">ログアウト</a></div></div><div class="user-menu-list" ng-hide="$root.User.loggedIn()"><div><a class="btn mypage-btn" ui-sref="register" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_register\'}">会員登録</a></div><div><a class="btn logout-btn" ui-sref="login" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_login\'}">ログイン</a></div></div><div ng-if="nb.notifications.length > 0"><div class="news-list-title">最新情報</div><ul class="news-list"><li class="menu-item menu-item-divided" ng-repeat="n in nb.notifications"><news data="n"></news></li><li class="menu-item menu-item-divided"><div class="more"><a ui-sref="notifications" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'notifications\'}">もっと見る</a></div></li></ul></div><ul class="other-menu-list"><li class="menu-item"><a ui-sref="faq" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_faq\'}"><span class="menu-text">よくある質問</span></a></li><li class="menu-item"><a ui-sref="contact" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_contact\'}"><span class="menu-text">お問い合わせ</span></a></li><li class="menu-item"><a ui-sref="policy" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_privacypolicy\'}"><span class="menu-text">プライバシーポリシー</span></a></li><li class="menu-item"><a ui-sref="links" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_links\'}"><span class="menu-text">当サイトへのリンクについて</span></a></li><li class="menu-item"><a ui-sref="terms" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_terms\'}"><span class="menu-text">利用規約</span></a></li><li class="menu-item"><a ui-sref="member_terms" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_member_terms\'}"><span class="menu-text">会員規約</span></a></li><li class="menu-item"><a ui-sref="company" class="menu-link" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_company\'}"><span class="menu-text">会社概要</span></a></li></ul><div class="footer"><div class="sns"><a class="facebook" href="#" socialshare="" socialshare-provider="facebook" socialshare-text="響 - HiBiKi Radio Station -" socialshare-url="https://hibiki-radio.jp" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_share_facebook\'}"><span class="icon"></span></a> <a class="twitter" ng-href="https://twitter.com/intent/follow?screen_name=HiBiKi_radio" target="_blank" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_follow_twitter\'}"><span class="icon"></span></a></div><div class="banners"><ul><li><a href="https://bushiroad.com/" target="_blank" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'banner_bushiroad\'}"><img src="/assets/images/bushiroad_banner.jpg" alt="ブシロード"></a></li><li><a href="https://www.bushiroad-music.com/" target="_blank" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'banner_bushiroad_music\'}"><img src="/assets/images/bushiroad_music.jpg" alt="ブシロードMUSIC"></a></li><li><a href="https://www.bushimo.jp/" target="_blank" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'banner_bushiroad_media\'}"><img src="/assets/images/bushimo.jpg" alt="ブシモ"></a></li><li><a href="https://bushiroad-ecshop.com/" target="_blank" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'hibiki_ec_shop\'}"><img src="/assets/images/bushiroad_ec_shop_banner.jpeg" alt="ブシロード EC SHOP"></a></li></ul></div><div class="navbar-adsense-ad" ng-include="\'app/components/adsense/responsive-rectangle.html\'"></div><div class="jasrac"><div class="jasrac-inner"><div class="mark"><img src="/assets/images/jasrac.jpg" alt="JASRAC"></div><div class="jasrac-number"><div>JASRAC許諾番号</div><div>9017659001Y31016</div></div></div><div class="jasrac-inner"><div class="mark"><img src="/assets/images/elicense_logo.png" alt="JASRAC"></div><div class="jasrac-number" style="text-align:left;"><div style="visibility:hidden;">0000000000000000</div><div>ID36922</div><div style="visibility:hidden;">0000000000000000</div></div></div></div><div class="jasrac"></div><span class="copyright"><i class="fa fa-copyright"></i> bushiroad All Rights Reserved.</span></div></nav>'),
    $templateCache.put("app/components/navbar/navbar.sp.html", '<div id="menu"><div class="menu-cell"><daily-voice voice="nb.dailyvoice"></daily-voice></div><div class="menu-cell"><a class="menu-heading" href="/"><img src="/assets/images/brand-logo-sp.png" alt="HiBiKi Station"></a></div><div class="menu-cell"><div class="menu-btn" ng-class="{active:nb.toggleMenu}" ng-click="nb.toggleMenu = !nb.toggleMenu"><span class="rect1"></span> <span class="rect2"></span> <span class="rect3"></span></div></div></div><nav ng-class="{active: nb.toggleMenu}"><ul class="menu-list"><li class="menu-item"><a ui-sref="home" class="menu-link" ng-click="nb.menuTap()" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_home\'}"><span class="menu-text">トップ</span></a></li><li class="menu-item" ng-show="$root.User.loggedIn()"><a ui-sref="favorites" class="menu-link" ng-click="nb.menuTap()" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_user_favorite\'}"><span class="menu-text">お気に入り一覧</span></a></li><li class="menu-item"><a ui-sref="informations" class="menu-link" ng-click="nb.menuTap()" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_informations\'}"><span class="menu-text">お知らせ一覧</span></a></li><li class="menu-item user-menu"><div ng-if="!$root.User.loggedIn()"><a ui-sref="register" class="user-link half" ng-click="nb.menuTap()" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_register\'}">新規会員登録</a> <a ui-sref="login" class="user-link half" ng-click="nb.menuTap()" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_login\'}">ログイン</a></div><div ng-if="$root.User.loggedIn()"><a ui-sref="mypage" class="user-link half" ng-click="nb.menuTap()" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_mypage\'}">マイページ</a> <a ui-sref="logout" class="user-link half" ng-click="nb.menuTap()" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'menu_logout\'}">ログアウト</a></div></li></ul></nav><div class="menu-background" ng-class="{active: nb.toggleMenu}" ng-click="nb.toggleMenu = !nb.toggleMenu"></div>'),
    $templateCache.put("app/components/news/news.html", '<div class="news" ng-class="{large: news.large}"><a ng-href="{{news.data.link_url}}" target="_blank" ng-if="news.data.kind == 1 || news.data.kind == 2" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'notification\',label:news.data.id}"><div><span class="newslabel {{news.label_class}}">{{news.label}}</span><span class="publish-start-at">{{news.data.publish_start_at|moment:"YYYY/MM/DD(ddd)"}}</span></div><div class="news-title">{{news.data.name}}</div></a> <a ui-sref="informations.detail({info_id: news.data.link_url})" ng-if="news.data.kind == 3" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'notification\',label:news.data.id}"><div><span class="newslabel {{news.label_class}}">{{news.label}}</span><span class="publish-start-at">{{news.data.publish_start_at|moment:"YYYY/MM/DD(ddd)"}}</span></div><div class="news-title">{{news.data.name}}</div></a> <a ui-sref="program.detail({access_id: news.data.link_url})" ng-if="news.data.kind == 4" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'notification\',label:news.data.id}"><div><span class="newslabel {{news.label_class}}">{{news.label}}</span><span class="publish-start-at">{{news.data.publish_start_at|moment:"YYYY/MM/DD(ddd)"}}</span></div><div class="news-title">{{news.data.name}}</div></a></div>'),
    $templateCache.put("app/components/program-list/program-list.html", '<div id="program-list" ng-show="pl.loaded"><div class="tabset"><div class="title">曜日別番組</div><ul class="week"><li class=""><a href="#" ng-click="pl.showPrograms(\'mon\')" ng-class="{active: pl.isMonday()}" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_day_of_week\',label:1}">{{pl.dayOfWeekStr(\'mon\')}}</a></li><li class=""><a href="#" ng-click="pl.showPrograms(\'tue\')" ng-class="{active: pl.isTuesday()}" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_day_of_week\',label:2}">{{pl.dayOfWeekStr(\'tue\')}}</a></li><li class=""><a href="#" ng-click="pl.showPrograms(\'wed\')" ng-class="{active: pl.isWednesday()}" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_day_of_week\',label:3}">{{pl.dayOfWeekStr(\'wed\')}}</a></li><li class=""><a href="#" ng-click="pl.showPrograms(\'thu\')" ng-class="{active: pl.isThursday()}" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_day_of_week\',label:4}">{{pl.dayOfWeekStr(\'thu\')}}</a></li><li class=""><a href="#" ng-click="pl.showPrograms(\'fri\')" ng-class="{active: pl.isFriday()}" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_day_of_week\',label:5}">{{pl.dayOfWeekStr(\'fri\')}}</a></li><li class=""><a href="#" ng-click="pl.showPrograms(\'satsun\')" ng-class="{active: pl.isSaturdayOrSunday()}" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_day_of_week\',label:6}">{{pl.dayOfWeekStr(\'satsun\')}}</a></li><li class=""><a href="#" ng-click="pl.showPrograms(\'all\')" ng-class="{active: pl.isAllday()}" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'home_day_of_week\',label:7}">{{pl.dayOfWeekStr(\'all\')}}</a></li></ul></div><div class="programs" ng-cloak=""><ul infinite-scroll="pl.loadMore()" infinite-scroll-disabled="pl.busy" infinite-scroll-distance="3"><li ng-repeat="p in pl.selected_programs"><div ng-if="($index - 1) % 3 == 0 && $index < 10"><div class="program-box program-box-animate"><div class="program-box-inner"><div ng-if="!$root.platform.mobile" ng-include="\'app/components/adsense/custom.html\'"></div><div ng-if="$root.platform.mobile" ng-include="\'app/components/adsense/responsive-rectangle.html\'"></div><div class="episode-updated-at"><div class="icon-container"><span class="additional-video"><span class="small-font">PR</span></span></div><div class="updated-date-container"><div class="updated-date-container-inner"><span class="updated-date">&nbsp;</span></div></div></div></div></div></div><div ng-if="($index - 1) % 3 != 0 || $index >= 10"><a ui-sref="program.detail({access_id: p.access_id})" class="program-box program-box-animate" logicalog-on="click" logicalog-properties="{category:\'tap\', action: \'home_program\', label: p.id}"><div class="program-box-inner"><div class="thumbnail" ng-style="{\'background-image\': \'url({{p.sp_image_url}})\'}"></div><div class="title program-title-animate">{{p.name | truncate: pl.name_length }}</div><div class="latest-episode-name">{{p.latest_episode_name | truncate: pl.latest_episode_name_length }}</div><div class="cast">{{p.cast | truncate: pl.cast_length}}</div><div class="episode-updated-at"><div class="icon-container"><span class="new-program" ng-show="p.new_program_flg"><span class="small-font">新番組</span></span> <span class="additional-video" ng-show="p.additional_video_flg"><span class="small-font">楽屋裏</span></span> <span class="chapter" ng-show="p.chapter_flg"><span class="small-font">チャプター</span></span></div><div class="updated-date-container"><div class="updated-date-container-inner"><span class="new" ng-show="p.update_flg"><span class="small-font">NEW</span></span> <span class="updated-date" ng-if="p.episode_updated_at">{{p.episode_updated_at | moment:\'YYYY.MM.DD\'}}</span> <span class="updated-date" ng-if="!p.episode_updated_at">&nbsp;</span></div></div></div></div></a></div></li></ul><div class="text-center" ng-if="$root.platform.tablet"><a href="#" class="btn" ng-click="pl.loadMore()" ng-show="!pl.busy">さらに読み込む</a></div></div></div>'),
    $templateCache.put("app/components/program-player/program-player.html", '<div id="program-screen-container" ng-class="{active:pgp.player.screen_view && pgp.player.ctrl}" ng-cloak=""><div id="program-screen" class="" ng-class="{active:pgp.player.screen_view && pgp.player.ctrl}"><span class="close" ng-hide="pgp.player.fullscreen || pgp.player.chapter_view" ng-click="pgp.player.screen(false)"></span> <span class="display-chapter" ng-click="pgp.player.chapter_view = true" ng-if="pgp.player.current_chapter && pgp.player.current_chapter.name" ng-show="!pgp.player.chapter_view" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'display_chapter\',label: pgp.player.episode_id}"></span><div class="chapter-view" ng-if="pgp.player.current_chapter && pgp.player.current_chapter.name" ng-class="{active:pgp.player.chapter_view}"><div class="chapter-scroll-view"><div class="name">{{pgp.player.current_chapter.name}}</div><div class="description"><p>{{pgp.player.current_chapter.description}}</p></div></div><div class="toggle-btn-container" ng-click="pgp.player.chapter_view = false" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'close_chapter\',label:pgp.player.episode_id}"><span class="toggle-btn"></span></div></div></div></div><div id="player-ctrl" ng-class="{active:pgp.player.ctrl}" ng-cloak=""><div class="player-ctrl-container"><div class="player-ctrl-inner"><div class="btns"><div class="live-hide"><span class="prev-btn" ng-class="{active: pgp.player.hasPrev()}" ng-click="pgp.player.prev()"></span></div><div class="play-btn-cell"><div class="spinner" ng-show="pgp.player.loading"><div class="rect1"></div><div class="rect2"></div><div class="rect3"></div><div class="rect4"></div></div><span class="play-btn" ng-class="{play: !pgp.player.paused()}" ng-click="pgp.player.playToggle()" ng-show="!pgp.player.loading" logicalog-on="click" logicalog-properties="{category:\'tap\',action:pgp.player.paused() ? \'pause_episode\' : \'play_episode\',label: pgp.player.episode_id + (pgp.player.backstage ? \'-1\' : \'-0\')}"></span></div><div class="live live-show">LIVE</div><div class="live-hide"><span class="next-btn" ng-class="{active: pgp.player.hasNext()}" ng-click="pgp.player.next()"></span></div><div class="rewind-btn-cell live-hide"><i class="sec">10</i> <span class="rewind-btn" ng-click="pgp.player.rewind(10)" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'back10s_episode\',label: pgp.player.episode_id + (pgp.player.backstage ? \'-1\' : \'-0\')}"></span></div></div><div class="screen-ctrl"><div class="on-off"><div>スクリーン</div><div ng-click="pgp.player.screenToggle()"><span class="on" ng-class="{active: pgp.player.screen_view}">ON</span> / <span class="off" ng-class="{active: !pgp.player.screen_view}">OFF</span></div></div><div class="expand" ng-show="!pgp.player.fullscreen" ng-if="::pgp.platform.pc || pgp.platform.tablet"><span ng-click="pgp.player.fullScreen(true)"><i class="fa fa-expand fa-2x fa-fw"></i></span></div><div class="compress" ng-show="pgp.player.fullscreen" ng-if="::pgp.platform.pc || pgp.platform.tablet"><span ng-click="pgp.player.fullScreen(false)"><i class="fa fa-compress fa-2x fa-fw"></i></span></div><div class="fullscreen-btn" ng-click="pgp.player.fullscreenVideo()"><div class="icon"><span></span> <span></span></div></div></div><div class="volume" ng-if="::pgp.platform.pc"><div class="volume-icon"><i class="fa"></i></div><div class="volume-slider" ng-mousedown="pgp.player.volumeHandle($event)"><span class="volume-slider-background"></span> <span class="volume-slider-level"></span></div></div><div class="seek-bar"><div><span class="time"></span></div><div><div class="seek-container" ng-mousedown="pgp.player.seekHandle($event)"><span class="seek-background"></span> <span class="seek"></span><div class="chapter-point-container"><span class="chapter-point" ng-repeat="c in pgp.player.chapter" ng-style="pgp.player.chapterPos(c)" tooltip="" tooltip-class="player-tooltip" tooltip-side="top" tooltip-title="{{c.name}}" ng-class="{active: pgp.player.isChapterActive(c.start_time)}"></span></div></div></div><div><span class="duration"></span></div></div><div class="track marquee" tooltip="" tooltip-class="player-tooltip" tooltip-title="{{pgp.player.program_name}} - {{pgp.player.episode_name}}" ng-if="::pgp.platform.pc || pgp.platform.tablet"><p class="name"><span class="autoplay" ng-show="pgp.player.auto_play">[連続再生中]</span>{{pgp.player.program_name}} - {{pgp.player.episode_name}}</p></div></div></div></div><div id="fullscreen-mask" ng-show="pgp.player.fullscreen" ng-click="pgp.player.fullScreen(false)"></div>'),
    $templateCache.put("app/components/user/limited.html", '<div id="user-limited"><div class="title">ログインまたは会員登録が必要です</div><div><p>このコンテンツは会員限定コンテンツです。引き続きお楽しみ頂く為にはログインまたは会員登録が必要です。</p><p>会員登録をしていただくと、番組のお気入り登録、番組の楽屋裏の視聴、番組のチャプター表示ができます。</p></div><div class="text-center mt20"><a ui-sref="register" ng-click="closeThisDialog()" class="btn register-btn">会員登録</a> <a ui-sref="login" ng-click="closeThisDialog()" class="btn">ログイン</a></div></div>'),
    $templateCache.put("app/mypage/partials/_password_change.html", '<div id="password-change"><p>「パスワードを変更」ボタンを押すと変更が完了します。</p><p>変更完了後にログアウトされますので、変更したパスワードでログインし直してください。</p><br><form id="password-change-form" name="password_change_form" ng-submit="my.passwordChangeSubmit()"><div class="form-group"><label>現在のパスワード</label> <input id="" type="password" name="old_password" class="form-control" autocorrect="off" ng-model="my.password.old" placeholder="8文字以上の半角英数字を入力してください" ng-required="true" ng-minlength="8" ng-maxlength="32" ng-pattern="/^[0-9a-zA-Z-_\\!\\.]+$/"><div ng-messages="password_change_form.old_password.$error" ng-if="password_change_form.old_password.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p><p class="error-message" ng-message="pattern">半角英数字と記号（-_!.）で入力してください</p><p class="error-message" ng-message="minlength">8文字以上で入力してください</p><p class="error-message" ng-message="maxlength">32文字以下で入力してください</p></div></div><div class="form-group"><label>新しいパスワード</label> <input id="" type="password" name="new_password" class="form-control" autocorrect="off" ng-model="my.password.new" placeholder="8文字以上の半角英数字を入力してください" ng-required="true" ng-minlength="8" ng-maxlength="32" ng-pattern="/^[0-9a-zA-Z-_\\!\\.]+$/"><div ng-messages="password_change_form.new_password.$error" ng-if="password_change_form.new_password.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p><p class="error-message" ng-message="pattern">半角英数字と記号（-_!.）で入力してください</p><p class="error-message" ng-message="minlength">8文字以上で入力してください</p><p class="error-message" ng-message="maxlength">32文字以下で入力してください</p></div></div><div class="form-group"><label>新しいパスワード(確認)</label> <input id="" type="password" name="new_password_confirm" class="form-control" autocorrect="off" ng-model="my.password.new_confirm" placeholder="確認のため同じパスワードを入力してください" ng-required="true" compare-to="my.password.new"><div ng-messages="password_change_form.new_password_confirm.$error" ng-if="password_change_form.new_password_confirm.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p><p class="error-message" ng-message="compareTo">パスワードが一致していません</p></div></div><div class="btn-group"><button type="submit" class="btn" ng-disabled="password_change_form.$invalid">パスワードを変更</button></div></form></div>'),
    $templateCache.put("app/mypage/partials/_profile_edit.html", '<div id="profile-edit"><form id="edit-form" name="edit_form" ng-submit="my.editSubmit()"><div class="form-group"><label>メールアドレス</label> <input type="email" name="email" class="form-control" autocorrect="off" ng-model="my.user.email" placeholder="ログインに使用するメールアドレスを入力してください" ng-required="true"><div ng-messages="edit_form.email.$error" ng-if="edit_form.email.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p><p class="error-message" ng-message="email">メールアドレスを入力してください</p></div></div><div class="form-group"><label>ニックネーム</label> <input type="text" name="name" class="form-control" ng-model="my.user.name" placeholder="ログイン後に表示される名前を入力してください" ng-required="true"><div ng-messages="edit_form.name.$error" ng-if="edit_form.name.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p></div></div><div class="row"><div class="col"><div class="form-group"><label>都道府県</label> <label for="prefecture" class="select-label"><select id="prefecture" class="form-control" name="prefecture" ng-model="my.user.prefecture_id" ng-options="option.value as option.label group by option.group for option in my.prefectures" ng-required="true"><option value="">都道府県を選択してください</option></select></label><div ng-messages="edit_form.prefecture.$error" ng-if="edit_form.prefecture.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p></div></div></div><div class="col"><div class="form-group"><label>生年月日</label><div class="birth-date"><label for="birth-date-year" class="select-label"><select id="birth-date-year" class="form-control" name="birth_date[year]" ng-model="my.user.birth_date_parts.year" ng-options="option for option in my.birth_date.year" ng-required="true"><option value=""></option></select></label> <span>年</span> <label for="birth-date-month" class="select-label"><select id="birth-date-month" class="form-control" name="birth_date[month]" ng-model="my.user.birth_date_parts.month" ng-options="option for option in my.birth_date.month" ng-required="true"><option value=""></option></select></label> <span>月</span> <label for="birth-date-day" class="select-label"><select id="birth-date-day" class="form-control" name="birth_date[day]" ng-model="my.user.birth_date_parts.day" ng-options="option for option in my.birth_date.day"><option value=""></option></select></label> <span>日</span></div><input type="hidden" name="birth_date" ng-model="my.user.birth_date" ng-required="true"><div ng-messages="edit_form.birth_date.$error" ng-if="edit_form.birth_date.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p></div></div></div></div><div class="row"><div class="col"><div class="form-group"><label>性別</label> <label><input id="" type="radio" name="gender" ng-model="my.user.gender" value="1"> 男</label> <label><input id="" type="radio" name="gender" ng-model="my.user.gender" value="2"> 女</label></div></div><div class="col"><div class="form-group allow-mail"><label>お知らせメールを</label> <label><input id="" type="radio" name="allow_mail_flg" ng-model="my.user.allow_mail_flg" value="true"> 受け取る</label> <label><input id="" type="radio" name="allow_mail_flg" ng-model="my.user.allow_mail_flg" value="false"> 受け取らない</label><div class="caption">響 -Hibiki Radio Station-からのオススメ情報やメールマガジンをお届けします</div><div ng-messages="edit_form.allow_mail_flg.$error" ng-if="edit_form.allow_mail_flg.$invalid" class="error-messages" role="alert"><p class="error-message" ng-message="required">入力は必須です</p></div></div></div></div><div class="btn-group"><button type="submit" class="btn" ng-disabled="edit_form.$invalid">更新</button></div></form></div>'),
    $templateCache.put("app/mypage/partials/_withdraw.html", '<div id="withdraw"><div ng-show="!my.withdraw_exec"><p>響 -HiBiKi Radio Station-を退会しますか？</p><p>退会をすると、会員限定の特典機能をお楽しみ頂く事が出来なくなります。</p><div class="mt30"><a class="btn red" href="#" ng-click="my.withdrawConfirm()">退会する</a></div></div><div class="nt30" ng-show="my.withdraw_exec"><p>響 Radio Stationをご利用頂きありがとうございました。</p><p>またのご利用をお待ちしております。</p><br><p><span>{{my.countdown}}</span>秒後にTOPページに移動します。</p></div></div>'),
    $templateCache.put("app/mypage/partials/_withdraw_confirm.html", '<p>響 -HiBiKi Radio Station-を退会しますか？</p><p>退会をすると、会員限定の特典機能をお楽しみ頂く事ができなくなります。</p><div class="mt30 text-center"><a class="btn" href="#" ng-click="closeThisDialog()">やっぱりやめる</a> <a class="btn red" href="#" ng-click="withdraw()">退会する</a></div>'),
    $templateCache.put("app/program/partials/_detail.html", '<div id="program-detail"><div class="sns"><ul><li><a class="facebook" href="#" socialshare="" socialshare-provider="facebook" socialshare-text="{{program.detail.share_text}}" socialshare-url="{{program.detail.share_url}}" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'tap_facebook\',label:program.detail.id}"><span class="icon"></span> Facebook</a></li><li><a class="twitter" href="#" socialshare="" socialshare-provider="twitter" socialshare-text="{{program.detail.share_text}}" socialshare-url="{{program.detail.share_url}}" socialshare-hashtags="{{program.twitter_hashtags}}" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'tap_twitter\',label:program.detail.id}"><span class="icon"></span> Twitter</a></li><li ng-if="!program.detail.user_favorite_flg"><a class="star" href="#" ng-click="program.addToFavorite($event)" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_register_user_favorite\',label:program.detail.id}"><span class="icon"></span> お気に入りに追加</a></li><li ng-if="program.detail.user_favorite_flg"><a class="unstar" href="#" ng-click="program.deleteToFavorite($event)" logicalog-properties="{category:\'tap\',action:\'program_unregister_user_favorite\',label:program.detail.id}"><span class="icon"></span> お気に入りから削除</a></li></ul></div><section class="detail program-box-animate" ng-animate-ref="program-box-{{program.access_id}}"><div class="box-inner-wrapper"><section class="box-inner-left"><h2 class="title program-title-animate" ng-bind="::program.detail.name" ng-animate-ref="program-title-{{program.access_id}}"></h2><p ng-bind-html="::program.detail.description | linky:\'_blank\'"></p><a class="official-site" ng-href="{{::program.detail.official_site_url}}" target="_blank" ng-if="::program.detail.official_site_url">公式サイトはこちら</a></section><section class="box-inner-right"><h2 class="title">ON AIR情報</h2><p ng-bind-html="::program.detail.onair_information | linky:\'_blank\'"></p></section></div></section><section class="personality" ng-if="::program.detail.casts.length > 0" ng-cloak=""><div class="title">パーソナリティ</div><div class="cast" ng-repeat="cast in ::program.detail.casts"><div class="thumbnail" ng-style="{\'background-image\': \'url({{cast.pc_image_url || \'/assets/images/no_image.png\'}})\'}" ng-if="::!program.isMobile" ng-cloak=""></div><div class="thumbnail" ng-style="{\'background-image\': \'url({{cast.sp_image_url || \'/assets/images/no_image.png\'}})\'}" ng-if="::program.isMobile"></div><div class="cast-name"><div class="name">{{::cast.name}}</div><div class="roll" ng-if="::cast.roll_name">{{::cast.roll_name}}</div></div></div></section><div class="episodes-wrapper" ng-if="::program.detail.episode.html_description || program.detail.episode.episode_parts.length > 0" ng-cloak=""><section class="episodes"><h2 class="title" ng-bind="::program.detail.episode.name"></h2><section ng-repeat="p in ::program.detail.episode.episode_parts" ng-if="::!program.detail.episode.html_description"><div class="episode-image" ng-if="::!program.isMobile && p.pc_image_url"><img ng-src="{{::p.pc_image_url}}"></div><div class="episode-image" ng-if="::program.isMobile && p.sp_image_url"><img ng-src="{{::p.sp_image_url}}"></div><p ng-bind-html="::p.description | linky:\'_blank\'"></p></section><section ng-bind-html="::program.detail.episode.html_description" ng-if="::program.detail.episode.html_description"></section></section></div><section class="links" ng-if="::program.detail.program_links.length > 0" ng-cloak=""><div class="title">関連リンク</div><div class="program-links"><ul><li ng-repeat="link in ::program.detail.program_links"><a ng-href="{{::link.link_url}}" target="_blank" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_relate_link\',label:link.id}"><img ng-src="{{::link.pc_image_url}}" alt="{{::link.name}}" ng-if="::link.pc_image_url"><p ng-if="::!link.pc_image_url">{{::link.name}}</p></a></li></ul></div></section></div>'),
    $templateCache.put("app/program/partials/_informations.html", '<div id="program-informations"><section class="nocontent-box" ng-hide="program.informations.length > 0"><p>お知らせはありません</p></section><section class="info" ng-repeat="info in program.informations"><span class="day">{{info.day}}</span><h2 class="title">{{info.name}}</h2><div ng-if="info.html_description"><div ng-bind-html="info.html_description"></div><div class="link" ng-if="info.link_url"><a ng-href="{{info.link_url}}" target="_blank" ng-if="info.kind == 1 || info.kind == 2">詳しくはこちら</a> <a ui-sref="informations.detail({info_id: info.link_url})" ng-if="info.kind == 3">詳しくはこちら</a> <a ui-sref="program.detail({access_id: info.link_url})" ng-if="info.kind == 4">詳しくはこちら</a></div></div><div ng-if="!info.html_description"><div class="info-image" ng-if="!program.isMobile && info.pc_image_url"><img ng-src="{{info.pc_image_url}}" alt=""></div><div class="info-image" ng-if="program.isMobile && info.sp_image_url"><img ng-src="{{info.sp_image_url}}" alt=""></div><p ng-bind-html="info.description | linky:\'_blank\'"></p><div class="link" ng-if="info.link_url"><a ng-href="{{info.link_url}}" target="_blank" ng-if="info.kind == 1 || info.kind == 2">詳しくはこちら</a> <a ui-sref="informations.detail({info_id: info.link_url})" ng-if="info.kind == 3">詳しくはこちら</a> <a ui-sref="program.detail({access_id: info.link_url})" ng-if="info.kind == 4">詳しくはこちら</a></div></div></section></div>'),
    $templateCache.put("app/program/partials/_introductions.html", '<div id="program-introductions"><div class="put-letter-container" ng-if="program.detail.message_form_url"><a href="#" ng-click="program.letter(program.detail.message_form_url)" class="put-letter" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_message\',label:program.detail.id}"><i class="fa fa-file-text-o fa-lg"></i>お便りを出す</a><div id="mailto" ng-if="program.detail.email"><div class="caption">直接メールを送りたい方はこちら</div><a href="mailto:{{program.detail.email}}" target="_blank" class="put-letter" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_mail\',label:program.detail.id}"><i class="icon fa fa-envelope-o fa-lg"></i>{{program.detail.email}}</a></div></div><div class="introductions"><h2 class="title">コーナー紹介</h2><section class="nocontent" ng-hide="program.detail.segments.length > 0"><p>コーナー情報はありません</p></section><section class="corner" ng-repeat="s in program.detail.segments"><h3 class="name" ng-bind="s.name"></h3><div ng-bind-html="::s.html_description" ng-if="::s.html_description"></div><div ng-repeat="p in s.segment_parts" ng-if="::!s.html_description"><div class="img-container" ng-if="p.pc_image_url.length > 0"><img ng-src="{{::p.pc_image_url}}" alt="{{::s.name}}"></div><p ng-bind-html="::p.description | linky:\'_blank\'"></p></div></section></div><div class="put-letter-container" ng-if="program.detail.message_form_url"><a href="#" ng-click="program.letter(program.detail.message_form_url)" class="put-letter" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_message\',label:program.detail.id}"><i class="fa fa-file-text-o fa-lg"></i>お便りを出す</a><div id="mailto" ng-if="program.detail.email"><div class="caption">直接メールを送りたい方はこちら</div><a href="mailto:{{program.detail.email}}" target="_blank" class="put-letter" logicalog-on="click" logicalog-properties="{category:\'tap\',action:\'program_mail\',label:program.detail.id}"><i class="icon fa fa-envelope-o fa-lg"></i>{{program.detail.email}}</a></div></div></div>'),
    $templateCache.put("app/program/partials/_products.html", '<div id="program-products"><section class="nocontent-box" ng-hide="program.products.length > 0"><p>商品はありません</p></section><section class="product" ng-repeat="product in program.products"><span class="day">{{product.day}}</span><h2 class="title">{{product.name}}</h2><div class="html-description" ng-if="product.html_description"><div ng-bind-html="product.html_description"></div><div class="link" ng-if="product.link_url"><a ng-href="{{product.link_url}}" target="_blank"><img src="/assets/images/button_shop.png"></a></div></div><div ng-if="!product.html_description"><div class="product-description"><p ng-bind-html="product.description | linky:\'_blank\'"></p></div><div class="table"><div class="product-image"><div class="product-info"><img ng-src="{{product.pc_image_url}}" alt="" ng-if="!program.isMobile"> <img ng-src="{{product.sp_image_url}}" alt="" ng-if="program.isMobile"><div class="info"><div class="price"><span>価格</span>{{product.price}}</div><div class="release_date"><span>発売日</span>{{product.release_date}}</div><div class="selling_agency"><span>発売元／販売元</span>{{product.selling_agency}}</div></div></div><div class="bonus-info" ng-if="product.bonus_text"><div class="image"><img ng-src="{{product.bonus_pc_image_url}}" alt="" ng-if="!program.isMobile"> <img ng-src="{{product.bonus_sp_image_url}}" alt="" ng-if="program.isMobile"></div><div class="info"><div class="info-title">ブシロード EC SHOP 購入特典</div><p>{{product.bonus_text}}</p></div></div></div><div class="product-detail"><p ng-bind-html="product.content_description | linky:\'_blank\'"></p><div class="link" ng-if="product.link_url"><a ng-href="{{product.link_url}}" target="_blank"><img src="/assets/images/button_shop.png"></a></div></div></div></div></section></div>')
}
]);

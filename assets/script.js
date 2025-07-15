let ws = null;
        let username = 'User';
        let currentChannel = 'general';
        let channels = new Set(['general']);
        let isPageVisible = true;
        let titleBlinkInterval = null;
        let originalTitle = 'EchoRoom - Real-time Conversations';
        let unreadCount = 0;

        // Funny, informal sample usernames
        const sampleUsernames = [
            // Food & Drinks
            'CoffeeAddict42', 'PizzaLover99', 'TacoTuesday', 'BurritoWizard', 'DonutDestroyer',
            'JellybeanJedi', 'PancakePhantom', 'WaffleWarrior', 'CerealKiller', 'NachoAverage',
            'SpaghettiMonster', 'PickleRick', 'BaconBeast', 'CookieMonster', 'IceCreamDream',
            'MarshmallowMage', 'FoodieLife', 'ChefKiss', 'BakeOff', 'GrillMaster', 'SaladDays',
            'VeggieVibes', 'MeatLover', 'SweetTooth', 'SpiceItUp', 'FlavorTown', 'CherryBomb',
            'HotSauceHero', 'ChocolateChamp', 'VanillaVibes', 'StrawberrySmash', 'BlueberryBlast',
            'RamenRider', 'SushiSamurai', 'PastaPirate', 'BreadBandit', 'CheeseChaser',

            // Tech & Gaming
            'ZoomZoomer', 'KeyboardWarrior', 'CtrlAltDelete', 'WiFiWhisperer', 'ByteMe',
            '404NotFound', 'DebuggingDuck', 'JavaJunkie', 'HtmlHero', 'CssNinja',
            'GamerGuru', 'ConsoleCommander', 'PCMasterRace', 'MobilePlayer', 'RetroGamer',
            'PixelPirate', 'NPCWhisperer', 'LootGoblin', 'RespawnRanger', 'AchievementHunter',
            'CodeCrusher', 'BugSquasher', 'GitCommander', 'DataDragon', 'CloudSurfer',
            'BinaryBoss', 'AlgorithmAce', 'DatabaseDuke', 'ServerSage', 'NetworkNinja',
            'CyberSleuth', 'DigitalDetective', 'VirtualVoyager', 'QuantumQueen', 'TechTitan',

            // Internet Culture & Social Media
            'DiscordMod', 'RedditLurker', 'MemeDealer', 'TikTokTalent', 'InstagramInfluencer',
            'YoutubeYeller', 'TwitterTroll', 'SnapchatSage', 'LinkedInLegend', 'FacebookFan',
            'MemeQueen', 'ViralViper', 'TrendTracker', 'HashtagHero', 'InfluencerImposter',
            'ContentCreator', 'StreamingStarfish', 'PodcastPanda', 'BloggingBear', 'VlogVampire',

            // Animals & Creatures
            'LazyPanda', 'SleepyDragon', 'FluffyUnicorn', 'CrazyCat', 'HappyHippo',
            'SneakySquirrel', 'BouncyBunny', 'WiseOwl', 'PlayfulPenguin', 'CheekyChimp',
            'MajesticMoose', 'CuriousCrow', 'GentleGiraffe', 'BraveBear', 'CleverCoyote',
            'SwiftSwallow', 'ProudPeacock', 'ShyShark', 'BoldBadger', 'KindKoala',
            'ElephantEars', 'ZebraStripes', 'LionHeart', 'TigerTail', 'WolfPack',

            // Personality & Vibes
            'ChaosMaster', 'SockThief', 'ProcrastinatorPro', 'WifiStealer', 'CouchPotato',
            'SnackAttack', 'DancingTaco', 'GiggleBox', 'ChillVibes', 'GoodVibesOnly',
            'StayPositive', 'SmileMore', 'LaughOutLoud', 'HappyFeet', 'SunnyDays',
            'RainbowChaser', 'StarGazer', 'DreamCatcher', 'AdventureSeeker', 'WanderlustWarrior',
            'FreeSpiritFly', 'ZenMaster', 'PeacefulPanda', 'JoyfulJester', 'OptimisticOwl',
            'EnergeticElf', 'CalmCactus', 'RelaxedRaven', 'SereneSeahorse', 'TranquilTurtle',

            // Hobbies & Activities
            'BookWorm', 'MovieBuff', 'MusicMaven', 'ArtAttack', 'DanceFloor',
            'SingASong', 'PaintBrush', 'CraftyCrafter', 'DIYExpert', 'HandyHelper',
            'GymRat', 'YogaMaster', 'RunnerHigh', 'SwimFast', 'BikeRider',
            'HikerLife', 'CampingFan', 'NatureLover', 'BeachBum', 'MountainMan',
            'PhotographyPhantom', 'GardeningGuru', 'CookingCaptain', 'TravelTiger', 'ReadingRaccoon',
            'WritingWizard', 'DrawingDragon', 'SingingSloth', 'DancingDolphin', 'YogaYak',

            // Random Fun
            'NoodleNinja', 'BubblePopper', 'CloudWatcher', 'PuddleJumper', 'StarCollector',
            'MoonDancer', 'SunChaser', 'RainMaker', 'WindWhisperer', 'ThunderClap',
            'LightningBolt', 'SnowFlake', 'IceCube', 'FireCracker', 'SparkPlug',
            'RocketFuel', 'TurboBoost', 'NitroNinja', 'SpeedDemon', 'FastTrack',
            'QuickSilver', 'RapidFire', 'BlazingBolt', 'SwiftStrike', 'FlashDash',
            'ZoomZap', 'WhirlWind', 'TornadoTwist', 'CycloneSpinner', 'HurricaneHunter',

            // Silly & Absurd
            'ToasterStrudel', 'WackyWombat', 'SillyGoose', 'CrazyCarrot', 'FunkyFerret',
            'WobbleWalrus', 'JumpingJellybean', 'BouncingBanana', 'TwirlingTurnip', 'SpinningSpoon',
            'FlippingFlapjack', 'RollingDonut', 'SlidingSloth', 'SkippingStone', 'HoppingHamster',
            'GlidingGopher', 'SoaringSquid', 'FloatingFish', 'DriftingDuck', 'CruisingCrab',
            'WaddlingWalrus', 'StrollingStarfish', 'Ambling Aardvark', 'MeanderingMole', 'RoamingRhino'
        ];

        function getRandomUsername() {
            const randomIndex = Math.floor(Math.random() * sampleUsernames.length);
            return sampleUsernames[randomIndex];
        }

        // Time formatting utility functions
        function getRelativeTime(timestamp) {
            try {
                const now = new Date();
                const messageTime = new Date(timestamp);
                
                // Check if the date is valid
                if (isNaN(messageTime.getTime())) {
                    console.log('Invalid timestamp:', timestamp);
                    return 'just now';
                }
                
                const diffInSeconds = Math.floor((now - messageTime) / 1000);

                if (diffInSeconds < 60) {
                    return diffInSeconds <= 10 ? 'just now' : `${diffInSeconds}s ago`;
                }

                const diffInMinutes = Math.floor(diffInSeconds / 60);
                if (diffInMinutes < 60) {
                    return diffInMinutes === 1 ? '1 min ago' : `${diffInMinutes} min ago`;
                }

                const diffInHours = Math.floor(diffInMinutes / 60);
                if (diffInHours < 24) {
                    return diffInHours === 1 ? '1 hour ago' : `${diffInHours} hours ago`;
                }

                const diffInDays = Math.floor(diffInHours / 24);
                if (diffInDays < 30) {
                    return diffInDays === 1 ? '1 day ago' : `${diffInDays} days ago`;
                }

                // For older messages, show date
                return messageTime.toLocaleDateString([], { 
                    month: 'short', 
                    day: 'numeric',
                    year: messageTime.getFullYear() !== now.getFullYear() ? 'numeric' : undefined
                });
            } catch (e) {
                console.error('Error in getRelativeTime:', e, 'for timestamp:', timestamp);
                return 'just now';
            }
        }

        // Notification and title management functions
        function updateNotificationStatus() {
            const statusEl = document.getElementById('notificationStatus');
            const iconEl = document.getElementById('notificationIcon');
            const textEl = document.getElementById('notificationText');
            
            if ('Notification' in window) {
                const permission = Notification.permission;
                statusEl.style.display = 'flex';
                statusEl.className = 'notification-status';
                
                if (permission === 'granted') {
                    statusEl.classList.add('enabled');
                    iconEl.textContent = '🔔';
                    textEl.textContent = 'Notifications enabled';
                } else if (permission === 'denied') {
                    statusEl.classList.add('denied');
                    iconEl.textContent = '🔕';
                    textEl.textContent = 'Notifications blocked';
                } else {
                    iconEl.textContent = '🔔';
                    textEl.textContent = 'Click to enable notifications';
                    statusEl.style.cursor = 'pointer';
                    statusEl.onclick = requestNotificationPermission;
                }
            } else {
                statusEl.style.display = 'none';
            }
        }
        
        function requestNotificationPermission() {
            if ('Notification' in window && Notification.permission === 'default') {
                Notification.requestPermission().then(permission => {
                    console.log('Notification permission:', permission);
                    updateNotificationStatus();
                });
            }
        }

        function playNotificationSound() {
            try {
                // Create a simple notification beep using Web Audio API
                const audioContext = new (window.AudioContext || window.webkitAudioContext)();
                const oscillator = audioContext.createOscillator();
                const gainNode = audioContext.createGain();
                
                oscillator.connect(gainNode);
                gainNode.connect(audioContext.destination);
                
                oscillator.frequency.setValueAtTime(800, audioContext.currentTime);
                oscillator.frequency.setValueAtTime(600, audioContext.currentTime + 0.1);
                
                gainNode.gain.setValueAtTime(0.1, audioContext.currentTime);
                gainNode.gain.exponentialRampToValueAtTime(0.01, audioContext.currentTime + 0.2);
                
                oscillator.start(audioContext.currentTime);
                oscillator.stop(audioContext.currentTime + 0.2);
            } catch (e) {
                console.log('Could not play notification sound:', e);
            }
        }

        function showNotification(title, body, icon = null) {
            if ('Notification' in window && Notification.permission === 'granted' && !isPageVisible) {
                const notification = new Notification(title, {
                    body: body,
                    icon: icon || 'data:image/svg+xml,%3Csvg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="%234a90e2"%3E%3Cpath d="M12 2C6.48 2 2 6.48 2 12c0 1.54.36 3.04 1.05 4.4L1 22l5.6-2.05C8.96 21.64 10.46 22 12 22c5.52 0 10-4.48 10-10S17.52 2 12 2zm-1 15h-2v-2h2v2zm0-4h-2V9h2v4z"%2F%3E%3C/svg%3E',
                    tag: 'echoroom-message',
                    requireInteraction: false
                });
                
                // Auto close after 5 seconds
                setTimeout(() => {
                    notification.close();
                }, 5000);
                
                // Focus window when notification is clicked
                notification.onclick = () => {
                    window.focus();
                    notification.close();
                };
                
                // Play sound only when notifications are enabled
                playNotificationSound();
            }
        }

        function startTitleBlink() {
            if (titleBlinkInterval) return; // Already blinking
            
            let isOriginal = true;
            titleBlinkInterval = setInterval(() => {
                if (isOriginal) {
                    document.title = `(${unreadCount}) New message${unreadCount > 1 ? 's' : ''} - EchoRoom`;
                } else {
                    document.title = originalTitle;
                }
                isOriginal = !isOriginal;
            }, 1500); // Slow blink every 1.5 seconds
        }

        function stopTitleBlink() {
            if (titleBlinkInterval) {
                clearInterval(titleBlinkInterval);
                titleBlinkInterval = null;
                document.title = originalTitle;
                unreadCount = 0;
            }
        }

        // Page visibility API to track when user is active
        document.addEventListener('visibilitychange', () => {
            isPageVisible = !document.hidden;
            if (isPageVisible) {
                stopTitleBlink();
            }
        });

        // Track window focus
        window.addEventListener('focus', () => {
            isPageVisible = true;
            stopTitleBlink();
        });

        window.addEventListener('blur', () => {
            isPageVisible = false;
        });

        function showLoadingSpinner(show = true) {
            const status = document.getElementById('status');
            const spinner = document.getElementById('loadingSpinner');
            const statusContent = status.querySelector('.status-content');
            
            if (show) {
                status.classList.add('connecting');
                spinner.style.display = 'flex';
                statusContent.style.display = 'none';
            } else {
                status.classList.remove('connecting');
                spinner.style.display = 'none';
                statusContent.style.display = 'flex';
            }
        }

        function initializeRandomUsername() {
            const randomUsername = getRandomUsername();
            username = randomUsername;
            document.getElementById('username').value = randomUsername;
        }

        function connect() {
            const status = document.getElementById('status');
            const messageInput = document.getElementById('messageInput');
            const sendButton = document.getElementById('sendButton');

            // Show loading spinner while connecting
            showLoadingSpinner(true);

            ws = new WebSocket('ws://localhost:8080/ws');

            ws.onopen = function () {
                showLoadingSpinner(false);
                const statusContent = status.querySelector('.status-content');
                const statusText = statusContent.querySelector('.status-text');
                statusText.textContent = 'Connected';
                status.className = 'status connected';
                messageInput.disabled = false;
                sendButton.disabled = false;
                console.log('Connected to WebSocket');

                // Send username immediately upon connection
                const usernameMessage = {
                    username: username,
                    content: '',
                    type: 'user_connected',
                    channel: currentChannel
                };
                ws.send(JSON.stringify(usernameMessage));
            };

            ws.onmessage = function (event) {
                try {
                    const message = JSON.parse(event.data);
                    console.log('Received WebSocket message:', message);

                    if (message.type === 'channel_switch') {
                        currentChannel = message.channel;
                        updateCurrentChannelDisplay();
                        // Use setTimeout to ensure channel is added to UI first (in case of race condition with channel_created)
                        setTimeout(() => {
                            updateChannelActiveState(currentChannel);
                        }, 10);
                        clearMessages();
                        displayMessage(message);
                        return;
                    }

                    if (message.type === 'channel_created') {
                        const channelName = message.name;
                        const channelType = message.channel_type;
                        if (!channels.has(channelName)) {
                            channels.add(channelName);
                            addChannelToList(channelName, channelType);
                        }
                        return;
                    }

                    if (message.type === 'active_channels') {
                        updateActiveChannelsList(message.channels);
                        return;
                    }

                    if (message.type === 'channel_deleted') {
                        const channelName = message.content;
                        removeChannelFromList(channelName);

                        // If user is in the deleted channel, switch to general
                        if (currentChannel === channelName) {
                            switchChannel('general');
                        }
                        return;
                    }

                    displayMessage(message);
                    
                    // Show notification for new messages when page is not visible
                    if (message.type === 'message' && message.username !== username && !isPageVisible) {
                        if ('Notification' in window && Notification.permission === 'granted') {
                            unreadCount++;
                            showNotification(
                                `New message in #${currentChannel}`,
                                `${message.username}: ${message.content}`,
                            );
                            startTitleBlink();
                        }
                    }
                } catch (e) {
                    displayMessage({
                        username: 'System',
                        content: event.data,
                        type: 'message'
                    });
                }
            };

            ws.onclose = function () {
                showLoadingSpinner(false);
                const statusContent = status.querySelector('.status-content');
                const statusText = statusContent.querySelector('.status-text');
                statusText.textContent = 'Disconnected';
                status.className = 'status disconnected';
                messageInput.disabled = true;
                sendButton.disabled = true;
                console.log('Disconnected from WebSocket');

                // Show reconnecting message after 1 second
                setTimeout(() => {
                    if (!ws || ws.readyState === WebSocket.CLOSED) {
                        showLoadingSpinner(true);
                    }
                }, 1000);
                
                setTimeout(connect, 3000);
            };

            ws.onerror = function (error) {
                console.error('WebSocket error:', error);
            };
        }

        function sendMessage() {
            const messageInput = document.getElementById('messageInput');
            const message = messageInput.value.trim();

            if (message && ws && ws.readyState === WebSocket.OPEN) {
                const messageObj = {
                    username: username,
                    content: message,
                    type: 'message',
                    channel: currentChannel
                };

                ws.send(JSON.stringify(messageObj));
                messageInput.value = '';
            }
        }

        function switchChannel(channelName) {
            if (channelName === currentChannel) return;

            const messageObj = {
                type: 'join_channel',
                channel: channelName
            };

            ws.send(JSON.stringify(messageObj));

            updateChannelActiveState(channelName);
        }

        function createChannel() {
            console.log('createChannel function called');
            const input = document.getElementById('newChannelInput');
            const typeSelect = document.getElementById('channelTypeSelect');
            const channelName = input.value.trim();
            const channelType = typeSelect.value;

            console.log('Channel name:', channelName);
            console.log('Channel type:', channelType);
            console.log('WebSocket state:', ws ? ws.readyState : 'null');
            console.log('WebSocket OPEN constant:', WebSocket.OPEN);

            if (channelName && ws && ws.readyState === WebSocket.OPEN) {
                const createChannelMsg = {
                    type: 'create_channel',
                    name: channelName,
                    channel_type: channelType // Use 'channel_type' instead of 'type'
                };

                console.log('Sending message:', createChannelMsg);
                ws.send(JSON.stringify(createChannelMsg));
                input.value = '';
                console.log('Message sent, input cleared');
            } else {
                console.log('Channel creation blocked:');
                console.log('- Channel name exists:', !!channelName);
                console.log('- WebSocket exists:', !!ws);
                console.log('- WebSocket open:', ws && ws.readyState === WebSocket.OPEN);
            }
        }

        function addChannelToList(channelName, channelType) {
            const channelsList = document.getElementById('channelsList');
            const channelItem = document.createElement('li');
            channelItem.className = 'channel-item';
            channelItem.dataset.channel = channelName;
            channelItem.onclick = () => switchChannel(channelName);

            const nameSpan = document.createElement('span');
            nameSpan.className = 'channel-name';
            nameSpan.textContent = '#' + channelName;

            const typeSpan = document.createElement('span');
            typeSpan.className = `channel-type ${channelType}`;
            typeSpan.textContent = channelType === 'persistent' ? '💾' : '⚡';

            channelItem.appendChild(nameSpan);
            channelItem.appendChild(typeSpan);
            channelsList.appendChild(channelItem);
        }

        function updateChannelActiveState(channelName) {
            const channelItems = document.querySelectorAll('.channel-item');
            channelItems.forEach(item => {
                if (item.dataset.channel === channelName) {
                    item.classList.add('active');
                } else {
                    item.classList.remove('active');
                }
            });
        }

        function updateCurrentChannelDisplay() {
            const currentChannelDiv = document.getElementById('currentChannel');
            currentChannelDiv.textContent = `Current Channel: #${currentChannel}`;
        }

        function clearMessages() {
            const messagesDiv = document.getElementById('messages');
            messagesDiv.innerHTML = '';
        }

        function updateActiveChannelsList(activeChannels) {
            // Clear existing channels
            const channelsList = document.getElementById('channelsList');
            channelsList.innerHTML = '';

            // Reset channels set
            channels = new Set();

            // Add all active channels
            activeChannels.forEach(channelInfo => {
                if (typeof channelInfo === 'string') {
                    // Legacy format - treat as ephemeral
                    channels.add(channelInfo);
                    addChannelToList(channelInfo, 'ephemeral');
                } else {
                    // New format with type
                    channels.add(channelInfo.name);
                    addChannelToList(channelInfo.name, channelInfo.type);
                }
            });

            // Set the active channel
            updateChannelActiveState(currentChannel);
        }

        function removeChannelFromList(channelName) {
            // Remove from channels set
            channels.delete(channelName);

            // Remove from DOM
            const channelsList = document.getElementById('channelsList');
            const channelItems = channelsList.querySelectorAll('.channel-item');

            channelItems.forEach(item => {
                if (item.dataset.channel === channelName) {
                    item.remove();
                }
            });
        }

        function displayMessage(message) {
            const messagesDiv = document.getElementById('messages');
            const messageDiv = document.createElement('div');
            messageDiv.className = 'message';

            if (message.username === username) {
                messageDiv.classList.add('own');
            }

            if (message.type === 'channel_switch') {
                messageDiv.classList.add('notification');
                // Auto-remove channel switch messages after 3 seconds
                setTimeout(() => {
                    messageDiv.classList.add('fade-out');
                    setTimeout(() => {
                        if (messageDiv.parentNode) {
                            messageDiv.parentNode.removeChild(messageDiv);
                        }
                    }, 500); // Wait for fade-out animation to complete
                }, 3000);
            } else if (message.username === 'System' || message.type === 'system_message') {
                messageDiv.classList.add('system');
            }

            // Format timestamp (but don't show for system messages like channel switches)
            let timestampHtml = '';
            if (message.type !== 'channel_switch') {
                let timestampStr = '';
                let fullTimestamp = '';
                
                if (message.timestamp) {
                    // If message has timestamp from server (history messages)
                    fullTimestamp = message.timestamp;
                    try {
                        timestampStr = getRelativeTime(message.timestamp);
                    } catch (e) {
                        console.log('Error parsing timestamp:', message.timestamp);
                        timestampStr = 'unknown time';
                    }
                } else {
                    // If no timestamp, use current time (real-time messages)
                    fullTimestamp = new Date().toISOString();
                    timestampStr = 'just now';
                }
                
                timestampHtml = `<span class="timestamp" title="${new Date(fullTimestamp).toLocaleString()}" data-timestamp="${fullTimestamp}">${timestampStr}</span>`;
            }

            messageDiv.innerHTML = `
                ${timestampHtml}
                <span class="username">${message.username}:</span>
                <span class="content">${message.content}</span>
            `;

            messagesDiv.appendChild(messageDiv);
            messagesDiv.scrollTop = messagesDiv.scrollHeight;
        }

        // Event listeners
        document.getElementById('messageInput').addEventListener('keypress', function (e) {
            if (e.key === 'Enter') {
                sendMessage();
            }
        });

        document.getElementById('newChannelInput').addEventListener('keypress', function (e) {
            if (e.key === 'Enter') {
                createChannelWithSpinner();
            }
        });

        document.getElementById('username').addEventListener('change', function (e) {
            username = e.target.value || 'User';
        });

        // Channel click handlers
        document.getElementById('channelsList').addEventListener('click', function (e) {
            // Find the closest channel-item (handles clicks on child elements)
            const channelItem = e.target.closest('.channel-item');
            if (channelItem) {
                const channelName = channelItem.dataset.channel;
                switchChannel(channelName);
            }
        });

        function toggleTheme() {
            const body = document.body;
            const themeToggle = document.querySelector('.theme-toggle');

            body.classList.toggle('dark-theme');

            // Update toggle button icon
            if (body.classList.contains('dark-theme')) {
                themeToggle.textContent = '☀️';
                themeToggle.title = 'Toggle Light Mode';
                localStorage.setItem('theme', 'dark');
            } else {
                themeToggle.textContent = '🌙';
                themeToggle.title = 'Toggle Dark Mode';
                localStorage.setItem('theme', 'light');
            }
        }

        function initializeTheme() {
            const savedTheme = localStorage.getItem('theme');
            const body = document.body;
            const themeToggle = document.querySelector('.theme-toggle');

            if (savedTheme === 'dark') {
                body.classList.add('dark-theme');
                themeToggle.textContent = '☀️';
                themeToggle.title = 'Toggle Light Mode';
            } else {
                themeToggle.textContent = '🌙';
                themeToggle.title = 'Toggle Dark Mode';
            }
        }

        // Function to update all relative timestamps
        function updateTimestamps() {
            const timestamps = document.querySelectorAll('.timestamp[data-timestamp]');
            timestamps.forEach(timestampEl => {
                const originalTime = timestampEl.getAttribute('data-timestamp');
                if (originalTime) {
                    timestampEl.textContent = getRelativeTime(originalTime);
                }
            });
        }

        // Update timestamps every 5 seconds
        setInterval(updateTimestamps, 5000);

        // Add loading spinner for channel creation
        window.createChannelWithSpinner = function() {
            const button = document.querySelector('.channel-input button');
            const originalText = button.textContent;
            
            button.innerHTML = '<div class="spinner" style="width: 12px; height: 12px; border-width: 1px;"></div>';
            button.disabled = true;
            
            createChannel();
            
            // Reset button after a short delay
            setTimeout(() => {
                button.textContent = originalText;
                button.disabled = false;
            }, 1000);
        }

        // Initialize theme, random username and connect when page loads
        initializeTheme();
        initializeRandomUsername();
        updateNotificationStatus();
        requestNotificationPermission();
        connect();
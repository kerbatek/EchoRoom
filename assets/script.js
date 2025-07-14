let ws = null;
        let username = 'User';
        let currentChannel = 'general';
        let channels = new Set(['general']);

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

        function initializeRandomUsername() {
            const randomUsername = getRandomUsername();
            username = randomUsername;
            document.getElementById('username').value = randomUsername;
        }

        function connect() {
            const status = document.getElementById('status');
            const messageInput = document.getElementById('messageInput');
            const sendButton = document.getElementById('sendButton');

            ws = new WebSocket('ws://localhost:8080/ws');

            ws.onopen = function () {
                status.innerHTML = '<span class="status-dot"></span><span>Connected</span>';
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
                } catch (e) {
                    displayMessage({
                        username: 'System',
                        content: event.data,
                        type: 'message'
                    });
                }
            };

            ws.onclose = function () {
                status.innerHTML = '<span class="status-dot"></span><span>Disconnected</span>';
                status.className = 'status disconnected';
                messageInput.disabled = true;
                sendButton.disabled = true;
                console.log('Disconnected from WebSocket');

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
                if (message.timestamp) {
                    // If message has timestamp from server (history messages)
                    const date = new Date(message.timestamp);
                    timestampStr = date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
                } else {
                    // If no timestamp, use current time (real-time messages)
                    const now = new Date();
                    timestampStr = now.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
                }
                timestampHtml = `<span class="timestamp">${timestampStr}</span>`;
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
                createChannel();
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

        // Initialize theme, random username and connect when page loads
        initializeTheme();
        initializeRandomUsername();
        connect();
document.addEventListener('DOMContentLoaded', function () {
	// Элементы DOM
	const petImage = document.getElementById('pet-image');
	const toast = document.getElementById('toast');
	const coinBalanceEl = document.getElementById('coin-balance');
	const burgerBtn = document.getElementById('burger-btn');
	const burgerOverlay = document.getElementById('burger-menu-overlay');
	const closeBurgerBtn = burgerOverlay.querySelector('.close-burger');
	const shopBtn = document.getElementById('shop-btn');
	const shopModal = document.getElementById('shop-modal');
	const closeShopBtn = document.getElementById('close-shop');
	const shopTabs = document.querySelectorAll('.shop-tab');
	const shopItemsContainer = document.getElementById('shop-items');

	// Начальные данные
	let coinBalance = 100;
	let petStatus = {
		hunger: 80, // сытость
		energy: 60, // энергия
		mood: 70, // настроение
		hygiene: 90, // чистота
	};

	// Примеры товаров магазина для каждого раздела
	const foodItems = [
		{ id: 1, name: 'Корм премиум', price: 20, bonus: 10 },
		{ id: 2, name: 'Витаминный корм', price: 30, bonus: 15 },
	];
	const gameItems = [
		{ id: 3, name: 'Игровой набор', price: 25, bonus: 5 },
		{ id: 4, name: 'Пазл', price: 15, bonus: 3 },
	];
	const otherItems = [{ id: 5, name: 'Аксессуар', price: 10, bonus: 0 }];
	let currentShopTab = 'food';

	// Обновление баланса в хедере
	function updateBalanceDisplay() {
		coinBalanceEl.textContent = coinBalance + ' 💰';
	}

	// Загрузка случайного изображения питомца (пример с Picsum)
	function loadPetImage() {
		const seed = Math.floor(Math.random() * 1000);
		petImage.src = `https://picsum.photos/seed/${seed}/300`;
	}

	// Функция для уведомлений (toast)
	function showToast(message) {
		toast.textContent = message;
		toast.classList.add('show');
		setTimeout(() => {
			toast.classList.remove('show');
		}, 2000);
	}

	// Обновление прогресс-баров для состояния питомца
	function updateStatusBars() {
		Object.keys(petStatus).forEach((key) => {
			const progressFill = document.querySelector(`.progress-bar[data-status="${key}"] .progress-fill`);
			if (progressFill) {
				progressFill.style.width = petStatus[key] + '%';
				progressFill.style.backgroundColor = petStatus[key] < 30 ? '#d0021b' : '#4a90e2';
			}
		});
	}

	// Обработка бургер-меню
	burgerBtn.addEventListener('click', () => {
		burgerOverlay.style.display = 'flex';
	});
	closeBurgerBtn.addEventListener('click', () => {
		burgerOverlay.style.display = 'none';
	});
	// Обработка кликов по пунктам меню
	burgerOverlay.querySelectorAll('a').forEach((link) => {
		link.addEventListener('click', function (e) {
			e.preventDefault();
			showToast(`Выбрано: ${this.textContent}`);
			burgerOverlay.style.display = 'none';
		});
	});

	// Открытие/закрытие магазина
	shopBtn.addEventListener('click', () => {
		shopModal.style.display = 'flex';
	});
	closeShopBtn.addEventListener('click', () => {
		shopModal.style.display = 'none';
	});

	// Переключение табов в магазине
	shopTabs.forEach((tab) => {
		tab.addEventListener('click', function () {
			shopTabs.forEach((t) => t.classList.remove('active'));
			this.classList.add('active');
			currentShopTab = this.dataset.tab;
			renderShopItems();
		});
	});

	// Отрисовка товаров магазина
	function renderShopItems() {
		shopItemsContainer.innerHTML = '';
		let items = [];
		if (currentShopTab === 'food') {
			items = foodItems;
		} else if (currentShopTab === 'games') {
			items = gameItems;
		} else if (currentShopTab === 'other') {
			items = otherItems;
		}
		items.forEach((item) => {
			const div = document.createElement('div');
			div.classList.add('shop-item');
			div.innerHTML = `
        <div class="item-name">${item.name}</div>
        <div class="item-price">${item.price} 💰</div>
      `;
			// При клике происходит проверка баланса и "покупка"
			div.addEventListener('click', () => {
				if (coinBalance >= item.price) {
					coinBalance -= item.price;
					updateBalanceDisplay();
					showToast(`Куплено: ${item.name}`);
					// Здесь можно добавить логику применения бонуса или выбора товара
				} else {
					showToast('Недостаточно монет');
				}
			});
			shopItemsContainer.appendChild(div);
		});
	}

	// Обработка кнопок действий с питомцем
	document.querySelectorAll('.action-btn').forEach((btn) => {
		btn.addEventListener('click', function () {
			const action = this.dataset.action;
			switch (action) {
				case 'feed':
					petStatus.hunger = Math.min(petStatus.hunger + 15, 100);
					showToast('Питомца покормили');
					break;
				case 'play':
					petStatus.mood = Math.min(petStatus.mood + 10, 100);
					petStatus.energy = Math.max(petStatus.energy - 10, 0);
					showToast('С питомцем поиграли');
					break;
				case 'sleep':
					petStatus.energy = Math.min(petStatus.energy + 20, 100);
					showToast('Питомец спит');
					break;
				case 'care':
					petStatus.hygiene = Math.min(petStatus.hygiene + 20, 100);
					showToast('За питомцем ухаживали');
					break;
			}
			// Обновление изображения питомца и его состояний
			loadPetImage();
			updateStatusBars();
		});
	});

	// Инициализация
	updateBalanceDisplay();
	loadPetImage();
	updateStatusBars();
	renderShopItems();
});

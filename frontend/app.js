// ═══════════════════ VIEWS LOADER ═══════════════════
async function loadViews() {
  const views = ['auth', 'dashboard', 'portfolio'];
  const htmlParts = await Promise.all(
    views.map(v => fetch(`views/${v}.html`).then(r => r.text()))
  );
  document.getElementById('app').innerHTML = htmlParts.join('');
}

// ═══════════════════ STATE ═══════════════════
let currentUser = null;
let currentPf   = null;

// ═══════════════════ VIEWS ═══════════════════
function showView(name) {
  document.querySelectorAll('.view').forEach(v => v.classList.remove('active'));
  document.getElementById('view-' + name).classList.add('active');
  if (name === 'dashboard') loadDashboard();
  if (name === 'portfolio' && currentPf) { renderPortfolio(); loadTransactions(); }
}

// ═══════════════════ AUTH TABS ═══════════════════
function switchTab(name, el) {
  document.querySelectorAll('.tab').forEach(t => t.classList.remove('active'));
  el.classList.add('active');
  document.getElementById('tab-login').style.display    = name === 'login'    ? '' : 'none';
  document.getElementById('tab-register').style.display = name === 'register' ? '' : 'none';
  clearMsg('msg-auth');
}

// ═══════════════════ API ═══════════════════
async function api(method, path, body) {
  const opts = { method, credentials: 'include', headers: {} };
  if (body) {
    opts.headers['Content-Type'] = 'application/json';
    opts.body = JSON.stringify(body);
  }
  const res  = await fetch(path, opts);
  const text = await res.text();
  if (!res.ok) throw new Error(text.trim() || 'HTTP ' + res.status);
  try { return JSON.parse(text); } catch { return text.trim(); }
}

// ═══════════════════ TOAST ═══════════════════
let toastTimer;
function toast(msg, type = 'ok') {
  const el = document.getElementById('toast');
  el.textContent = (type === 'ok' ? '✓ ' : '✗ ') + msg;
  el.className = 'toast ' + type;
  el.style.display = 'block';
  clearTimeout(toastTimer);
  toastTimer = setTimeout(() => el.style.display = 'none', 3500);
}

function showMsg(id, text, type = 'error') {
  const el = document.getElementById(id);
  el.textContent = (type === 'error' ? '✗ ' : '✓ ') + text;
  el.className = 'msg msg-' + (type === 'error' ? 'error' : 'success') + ' show';
}
function clearMsg(id) {
  const el = document.getElementById(id);
  if (el) el.className = 'msg';
}

// ═══════════════════ SIGN IN ═══════════════════
async function handleSignIn(e) {
  e.preventDefault();
  clearMsg('msg-auth');
  try {
    const username = await api('POST', '/signin', {
      email:    document.getElementById('l-email').value,
      password: document.getElementById('l-pass').value,
    });
    await loadUser(username);
    showView('dashboard');
  } catch {
    showMsg('msg-auth', 'INVALID EMAIL OR PASSWORD');
  }
}

// ═══════════════════ SIGN UP ═══════════════════
async function handleSignUp(e) {
  e.preventDefault();
  clearMsg('msg-auth');
  try {
    await api('POST', '/signup', {
      username: document.getElementById('r-user').value,
      email:    document.getElementById('r-email').value,
      password: document.getElementById('r-pass').value,
    });
    showMsg('msg-auth', 'ACCOUNT CREATED — PLEASE SIGN IN', 'success');
    switchTab('login', document.getElementById('tab-login-btn'));
  } catch (err) {
    showMsg('msg-auth', err.message.toUpperCase());
  }
}

// ═══════════════════ LOGOUT ═══════════════════
function handleLogout() {
  document.cookie = 'X-Auth=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;';
  localStorage.removeItem('tradeexe_user');
  currentUser = null;
  currentPf   = null;
  showView('auth');
}

// ═══════════════════ USER ═══════════════════
async function loadUser(username) {
  const data = await api('GET', '/users/' + username);
  currentUser = data;
  document.getElementById('h-user').textContent   = data.username.toUpperCase();
  document.getElementById('h-user-2').textContent = data.username.toUpperCase();
  localStorage.setItem('tradeexe_user', data.username);
}

// ═══════════════════ DASHBOARD ═══════════════════
async function loadDashboard() {
  if (!currentUser) return;
  const listEl = document.getElementById('portfolio-list');
  listEl.innerHTML = '<div class="loading">LOADING <span class="blink">█</span></div>';
  try {
    const data = await api('GET', '/users/' + currentUser.username);
    currentUser = data;
    renderPortfolioList(data.portfolios || []);
  } catch {
    listEl.innerHTML = '<div class="loading" style="color:var(--red)">✗ FAILED TO LOAD</div>';
  }
}

function renderPortfolioList(portfolios) {
  const el = document.getElementById('portfolio-list');
  if (!portfolios.length) {
    el.innerHTML = `<div class="empty-state">▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓<br>NO PORTFOLIOS YET<br>CREATE ONE BELOW<br>▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓</div>`;
    return;
  }
  el.innerHTML = portfolios.map(p => `
    <div class="p-card" onclick='openPortfolio(${JSON.stringify(p)})'>
      <div>
        <div class="p-card-name">▶ ${p.name.toUpperCase()}</div>
        <div class="p-card-meta">COST: <span>${fmt(p.total_cost)}</span></div>
      </div>
      <div class="p-card-arrow">▶</div>
    </div>
  `).join('');
}

async function handleAddPortfolio(e) {
  e.preventDefault();
  clearMsg('msg-newportfolio');
  const name = document.getElementById('p-name').value.trim();
  try {
    await api('POST', `/users/${currentUser.username}/portfolio`, { name });
    document.getElementById('p-name').value = '';
    toast('PORTFOLIO CREATED: ' + name.toUpperCase());
    await loadDashboard();
  } catch (err) {
    showMsg('msg-newportfolio', err.message.toUpperCase());
  }
}

// ═══════════════════ PORTFOLIO ═══════════════════
function openPortfolio(p) {
  currentPf = p;
  showView('portfolio');
}

function renderPortfolio() {
  const p = currentPf;
  document.getElementById('pf-title').textContent = '◉ ' + p.name.toUpperCase();
  document.getElementById('s-cost').textContent   = fmt(p.total_cost);

  const pnl = parseFloat(p.total_profit_loss) || 0;
  const pnlEl = document.getElementById('s-pnl');
  pnlEl.textContent = (pnl >= 0 ? '+' : '') + fmt(p.total_profit_loss);
  pnlEl.className = 'stat-val ' + (pnl > 0 ? 'green' : pnl < 0 ? 'red' : 'blue');

  document.getElementById('s-created').textContent =
    p.created ? new Date(p.created).toLocaleDateString('pl-PL') : '—';
}

async function handleAddTransaction(e) {
  e.preventDefault();
  clearMsg('msg-tx');
  const sym   = document.getElementById('tx-sym').value.trim();
  const qty   = document.getElementById('tx-qty').value;
  const price = document.getElementById('tx-price').value;
  const type  = document.getElementById('tx-type').value;

  try {
    await api('POST', `/users/${currentUser.username}/portfolio/${currentPf.name}/transactions`, {
      symbol: sym, quantity: qty, amount: price, type,
    });

    document.getElementById('tx-sym').value   = '';
    document.getElementById('tx-qty').value   = '';
    document.getElementById('tx-price').value = '';
    document.getElementById('order-preview').style.display = 'none';

    toast(`${type.toUpperCase()} ${qty} ${sym} @ $${parseFloat(price).toFixed(2)}`);

    const data = await api('GET', '/users/' + currentUser.username);
    currentUser = data;
    const updated = (data.portfolios || []).find(p => p.name === currentPf.name);
    if (updated) { currentPf = updated; renderPortfolio(); loadTransactions(); }

  } catch (err) {
    showMsg('msg-tx', err.message.toUpperCase());
  }
}

// ═══════════════════ TRANSACTIONS ═══════════════════
async function loadTransactions() {
  const el = document.getElementById('tx-list');
  try {
    const data = await api('GET', `/users/${currentUser.username}/portfolio/${currentPf.name}/transactions`);
    const txns = data.transactions || [];
    if (!txns.length) {
      el.innerHTML = '<div class="empty-state">NO TRANSACTIONS YET</div>';
      return;
    }
    el.innerHTML = `
      <table class="tx-table">
        <thead><tr>
          <th>TYPE</th><th>SYMBOL</th><th>QTY</th><th>PRICE</th><th>TOTAL</th><th>DATE</th>
        </tr></thead>
        <tbody>${txns.map(t => {
          const isBuy = String(t.type) === '0';
          const total = (parseFloat(t.quantity) * parseFloat(t.price)).toFixed(2);
          return `<tr>
            <td class="${isBuy ? 'tx-buy' : 'tx-sell'}">${isBuy ? '▲ BUY' : '▼ SELL'}</td>
            <td style="color:var(--blue)">${t.symbol}</td>
            <td>${parseFloat(t.quantity)}</td>
            <td>${fmt(t.price)}</td>
            <td style="color:var(--gold)">$${total}</td>
            <td>${new Date(t.created).toLocaleDateString('pl-PL')}</td>
          </tr>`;
        }).join('')}</tbody>
      </table>`;
  } catch {
    el.innerHTML = '<div class="loading" style="color:var(--red)">✗ FAILED TO LOAD</div>';
  }
}

// ═══════════════════ ORDER PREVIEW ═══════════════════
function updatePreview() {
  const qty   = parseFloat(document.getElementById('tx-qty').value)   || 0;
  const price = parseFloat(document.getElementById('tx-price').value) || 0;
  const type  = document.getElementById('tx-type').value;
  const sym   = document.getElementById('tx-sym').value || '—';
  const el    = document.getElementById('order-preview');

  if (qty > 0 && price > 0) {
    document.getElementById('order-total').textContent = fmt(qty * price);
    document.getElementById('order-desc').textContent  =
      `${type.toUpperCase()} ${qty} × ${sym} @ $${price}`;
    el.style.display = 'block';
  } else {
    el.style.display = 'none';
  }
}

// ═══════════════════ HELPERS ═══════════════════
function fmt(val) {
  const n = parseFloat(val) || 0;
  return '$' + n.toFixed(2);
}

// ═══════════════════ INIT ═══════════════════
(async () => {
  await loadViews();
  const saved = localStorage.getItem('tradeexe_user');
  if (saved) {
    try {
      await loadUser(saved);
      showView('dashboard');
    } catch {
      localStorage.removeItem('tradeexe_user');
    }
  }
})();

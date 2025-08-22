

function formatCurrency(amount, currency = 'USD') {
  const validCurrencies = new Set([
    'USD', 'EUR', 'RUB', 'GBP', 'JPY', 'CNY', 'KZT', 'UAH', 'BYN', 'PLN', 'INR', 'BRL', 'CAD', 'AUD'
  ]);

  if (!currency || !/^[A-Z]{3}$/i.test(currency) || !validCurrencies.has(currency.toUpperCase())) {
    console.warn(`Invalid or unsupported currency: ${currency}. Using USD as fallback.`);
    currency = 'USD';
  }

  const formatted = new Intl.NumberFormat('en-US', {
    style: 'currency',
    currency: currency.toUpperCase()
  }).format(amount);

  return formatted;
}

function formatDate(dateString) {
  return new Date(dateString).toLocaleString('en-US', {
    year: 'numeric',
    month: 'long',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  });
}

function getStatusBadge(status) {
  const statusMap = {
    202: { class: 'status-202', text: 'Accepted' },
    200: { class: 'status-pending', text: 'Processing' },
    404: { class: 'status-pending', text: 'Pending' }
  };
  const statusInfo = statusMap[status] || { class: 'status-pending', text: `Status ${status}` };
  return `<span class="status-badge ${statusInfo.class}">${statusInfo.text}</span>`;
}

function displayOrder(orderData) {
  const orderDisplay = document.getElementById('orderDisplay');
  orderDisplay.innerHTML = `
    <div class="order-container">
      <div class="order-header">
        <div class="order-header-top">
          <div class="order-title-section">
            <h2 class="order-title">Order #${orderData.order_uid}</h2>
            <div class="order-subtitle">📋 Managed via ${orderData.delivery_service.toUpperCase()} delivery service</div>
          </div>
          <div class="order-status-section">
            <div class="order-date">Created ${formatDate(orderData.date_created)}</div>
            <div class="order-total-badge">${formatCurrency(orderData.payment.amount, orderData.payment.currency)}</div>
          </div>
        </div>
        <div class="order-meta">
          <div class="meta-item"><div class="meta-label">🚚 Track Number</div><div class="meta-value">${orderData.track_number}</div></div>
          <div class="meta-item"><div class="meta-label">👤 Customer ID</div><div class="meta-value">${orderData.customer_id}</div></div>
          <div class="meta-item"><div class="meta-label">📦 Delivery Service</div><div class="meta-value">${orderData.delivery_service.toUpperCase()}</div></div>
          <div class="meta-item"><div class="meta-label">💳 Transaction</div><div class="meta-value">${orderData.payment.transaction}</div></div>
        </div>
      </div>
      <div class="order-body">
        <div class="section-grid">
          <div class="section-card">
            <h3 class="section-title"><div class="section-icon">🚚</div> Delivery Information</h3>
            <div class="info-grid">
              <div class="info-row"><span class="info-label">Full Name</span><span class="info-value">${orderData.delivery.name}</span></div>
              <div class="info-row"><span class="info-label">Phone</span><span class="info-value">${orderData.delivery.phone}</span></div>
              <div class="info-row"><span class="info-label">Email</span><span class="info-value">${orderData.delivery.email}</span></div>
              <div class="info-row"><span class="info-label">Address</span><span class="info-value">${orderData.delivery.address}</span></div>
              <div class="info-row"><span class="info-label">City</span><span class="info-value">${orderData.delivery.city}</span></div>
              <div class="info-row"><span class="info-label">Region</span><span class="info-value">${orderData.delivery.region}</span></div>
              <div class="info-row"><span class="info-label">ZIP</span><span class="info-value">${orderData.delivery.zip}</span></div>
            </div>
          </div>
          <div class="section-card">
            <h3 class="section-title"><div class="section-icon">💳</div> Payment Breakdown</h3>
            <div class="info-grid">
              <div class="info-row"><span class="info-label">Provider</span><span class="info-value">${orderData.payment.provider.toUpperCase()}</span></div>
              <div class="info-row"><span class="info-label">Bank</span><span class="info-value">${orderData.payment.bank.toUpperCase()}</span></div>
              <div class="info-row"><span class="info-label">Items Subtotal</span><span class="info-value">${formatCurrency(orderData.payment.goods_total, orderData.payment.currency)}</span></div>
              <div class="info-row"><span class="info-label">Delivery Cost</span><span class="info-value">${formatCurrency(orderData.payment.delivery_cost, orderData.payment.currency)}</span></div>
              <div class="info-row"><span class="info-label">Custom Fees</span><span class="info-value">${formatCurrency(orderData.payment.custom_fee, orderData.payment.currency)}</span></div>
              <div class="info-row total-row"><span class="info-label">Total Amount</span><span class="info-value">${formatCurrency(orderData.payment.amount, orderData.payment.currency)}</span></div>
            </div>
          </div>
        </div>
        <div class="items-section">
          <div class="items-header">
            <h3 class="section-title"><div class="section-icon">📦</div> Order Items (${orderData.items.length})</h3>
          </div>
          <table class="items-table">
            <thead>
              <tr>
                <th>Product</th>
                <th>Brand</th>
                <th>ID</th>
                <th>Size</th>
                <th>Original</th>
                <th>Sale</th>
                <th>Total</th>
                <th>Status</th>
                <th>Track</th>
              </tr>
            </thead>
            <tbody>
              ${orderData.items.map(item => `
                <tr>
                  <td class="product-cell">${item.name}</td>
                  <td class="brand-cell">${item.brand}</td>
                  <td>${item.chrt_id}</td>
                  <td>${item.size || 'N/A'}</td>
                  <td class="price-cell"><span class="price-original">${formatCurrency(item.price, orderData.payment.currency)}</span></td>
                  <td><span class="price-sale">-${item.sale}%</span></td>
                  <td class="price-cell">${formatCurrency(item.total_price, orderData.payment.currency)}</td>
                  <td>${getStatusBadge(item.status)}</td>
                  <td><span class="track-cell">${item.track_number}</span></td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  `;
}

function showLoading() {
  document.getElementById('orderDisplay').innerHTML = `
    <div class="loading-state">
      <div class="loading-spinner"></div>
      <h3>Loading Order...</h3>
      <p>Please wait while we fetch the order details.</p>
    </div>
  `;
}

function showError(message) {
  document.getElementById('orderDisplay').innerHTML = `
    <div class="error-state">
      <div style="font-size: 48px; margin-bottom: 1rem;">😢</div>
      <h3>Order Not Found</h3>
      <p>${message}</p>
    </div>
  `;
}


async function searchOrder() {
  const uid = document.getElementById('orderIdInput').value.trim();
  console.log("Searching for order UID:", uid);

  if (!uid) {
    showError("Please enter an Order UID");
    return;
  }


  try {
    const response = await fetch(`/api/v1/order/${encodeURIComponent(uid)}`);
    console.log("Response status:", response.status); 

    if (!response.ok) {
      if (response.status === 404) {
        showError("The requested order does not exist.");
      } else {
        const errorText = await response.text();
        console.error("Error response:", errorText);
        showError(`Error ${response.status}: ${errorText || 'Unknown error'}`);
      }
      return;
    }

    const orderData = await response.json();
    console.log("Order data received:", orderData);
    displayOrder(orderData);

  } catch (error) {
    console.error("Fetch error:", error);
    showError("Failed to connect to the server");
  }
}


document.addEventListener('DOMContentLoaded', () => {
  document.getElementById('searchBtn').addEventListener('click', searchOrder);
  document.getElementById('orderIdInput').addEventListener('keypress', (e) => {
    if (e.key === 'Enter') searchOrder();
  });
});

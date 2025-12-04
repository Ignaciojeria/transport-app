import Link from 'next/link'
import Image from 'next/image'

export const metadata = {
  title: 'Refund Policy - MiCartaPro',
  description: 'MiCartaPro refund policy. Learn about our refund and cancellation terms.',
}

export default function RefundPage() {
  return (
    <div className="min-h-screen bg-gradient-to-br from-blue-50 via-white to-indigo-50">
      {/* Navigation */}
      <nav className="border-b bg-white/80 backdrop-blur-sm sticky top-0 z-50">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="flex justify-between items-center h-16">
            <Link href="/" className="flex items-center space-x-2">
              <Image 
                src="/logo.png" 
                alt="MiCartaPro Logo" 
                width={240} 
                height={72}
                className="h-16 md:h-20 w-auto"
              />
            </Link>
            <Link 
              href="/"
              className="text-gray-600 hover:text-blue-600 transition-colors"
            >
              Back to Home
            </Link>
          </div>
        </div>
      </nav>

      {/* Content */}
      <div className="max-w-4xl mx-auto px-4 sm:px-6 lg:px-8 py-12 md:py-16">
        <div className="bg-white rounded-lg shadow-lg p-8 md:p-12">
          <h1 className="text-3xl md:text-4xl font-bold text-gray-900 mb-6">
            Refund Policy
          </h1>
          
          <p className="text-gray-600 mb-8">
            Last updated: December 3, 2025
          </p>

          <div className="prose prose-lg max-w-none space-y-8 text-gray-700">
            <section>
              <p>
                MiCartaPro uses Paddle.com as our payment processor and Merchant of Record (MoR). Paddle handles payment processing, invoicing, taxes (where applicable), subscription renewals, and refund reviews under their consumer protection and legal framework.
              </p>
              <p>
                By completing a purchase, you agree to Paddle's Terms & Conditions and to this Refund & Subscription Policy.
              </p>
              <p className="mt-4">
                📄 <strong>Paddle Terms:</strong> <a href="https://paddle.com/legal" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:text-blue-700 underline">https://paddle.com/legal</a>
              </p>
              <p>
                📄 <strong>Paddle Privacy Policy:</strong> <a href="https://paddle.com/privacy" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:text-blue-700 underline">https://paddle.com/privacy</a>
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">1. Nature of the Product</h2>
              <p>
                MiCartaPro is a subscription-based digital software service (SaaS).
              </p>
              <p>
                Access to the dashboard and features is granted immediately after purchase, which constitutes delivery of digital content.
              </p>
              <p>
                Per Paddle's consumer terms:
              </p>
              <p className="italic">
                "Where a product is digital content which is immediately made available... you consent to immediate performance and acknowledge that you will lose your right of withdrawal once access begins."
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">2. 14-Day Refund Guarantee (Initial Purchase Only)</h2>
              <p>
                We offer a 14-day money-back guarantee for first-time subscription purchases only.
              </p>
              <p>
                If you are not satisfied, you may request a refund within the first 14 days — no explanation required.
              </p>
              <p>
                Refunds are processed exclusively by Paddle.
              </p>
              <p>
                Funds are returned to the original payment method.
              </p>
              <p>
                After 14 days, the purchase becomes final.
              </p>
              <p>
                ⏳ Refund requests submitted after 60 days may not be eligible, as defined by Paddle's policy.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">3. Refund Evaluation & Limitations</h2>
              <p>
                Refund approvals are at Paddle's discretion, based on their policies and applicable consumer protection laws.
              </p>
              <p>
                Refunds may be refused in cases of:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>evidence of service usage or benefit received</li>
                <li>refund abuse or manipulative behavior</li>
                <li>fraud or suspicious activity</li>
                <li>attempts to exploit the policy</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">4. Access Equals Delivery</h2>
              <p>
                Once the user logs in, accesses the dashboard, or uses the service, the product is considered delivered and consumed, and the statutory right of withdrawal may no longer apply.
              </p>
              <p>
                Examples that count as access:
              </p>
              <ul className="list-none pl-0 space-y-2">
                <li>✔ logging into the dashboard</li>
                <li>✔ creating or editing menu items</li>
                <li>✔ generating or using QR features</li>
                <li>✔ receiving onboarding or support</li>
                <li>✔ engaging with any core platform feature</li>
              </ul>
              <p className="mt-4">
                If you wish to request a refund, we recommend doing so before using the platform.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">5. Subscription Renewals & Cancellation</h2>
              <p>
                Subscriptions renew automatically.
              </p>
              <p>
                You may cancel at any time.
              </p>
              <p>
                To avoid being charged again, cancellation must be completed at least 48 hours before renewal.
              </p>
              <p>
                Cancellation prevents future charges — it does not trigger refunds for already-charged renewals.
              </p>
              <p>
                Refund eligibility applies only to the first billing, not to renewals or unused subscription time.
              </p>
              <p>
                No refunds are granted for:
              </p>
              <ul className="list-none pl-0 space-y-2">
                <li>✘ renewal payments already processed</li>
                <li>✘ unused time in a current billing period</li>
                <li>✘ change of mind after usage</li>
                <li>✘ lack of usage by customer choice</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">6. Add-On Human Services (Separate Billing)</h2>
              <p>
                We may offer optional add-on services such as:
              </p>
              <ul className="list-disc pl-6 space-y-2 mb-4">
                <li>menu setup / customization</li>
                <li>data import or migration</li>
                <li>product/catalog loading</li>
                <li>design, consulting or advanced support</li>
              </ul>
              <p>These services:</p>
              <ul className="list-disc pl-6 space-y-2">
                <li>are billed separately (Bank Transfer / PayPal / Invoice)</li>
                <li>are not refundable through Paddle</li>
                <li>follow their own independent service/refund terms</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">7. Billing, Taxes & Documentation</h2>
              <p>
                Paddle may apply local taxes (VAT/GST/consumption tax) depending on region.
              </p>
              <p>
                All invoices and receipts are delivered electronically.
              </p>
              <p>
                Customers are responsible for entering accurate billing information.
              </p>
              <p>
                If registered for indirect tax reclaim, you may request VAT/GST refund from Paddle within 60 days from purchase, subject to jurisdictional rules.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">8. Support & Refund Requests</h2>
              <p>
                <strong>MiCartaPro Support:</strong>
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>📧 support@micartapro.com</li>
                <li>📱 +56 9 5785 7558</li>
              </ul>
              <p className="mt-4">
                Refund or billing assistance via Paddle:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>🌐 <a href="https://paddle.net/support" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:text-blue-700 underline">https://paddle.net/support</a></li>
              </ul>
              <p className="mt-4">
                Paddle manages refund processing — MiCartaPro cannot issue refunds directly.
              </p>
            </section>
          </div>

          <div className="mt-12 pt-8 border-t border-gray-200 flex flex-col sm:flex-row justify-between items-center gap-4">
            <Link 
              href="/"
              className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
            >
              ← Back to Home
            </Link>
            <div className="flex flex-col sm:flex-row gap-4">
              <Link 
                href="/terms"
                className="inline-flex items-center text-blue-600 hover:text-blue-700 font-medium"
              >
                View Terms and Conditions →
              </Link>
            </div>
          </div>
        </div>
      </div>
    </div>
  )
}

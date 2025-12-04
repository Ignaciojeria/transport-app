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
                MiCartaPro uses Paddle.com as our payment processor and Merchant of Record (MoR).
              </p>
              <p>
                All payments, billing and refunds for purchases made through Paddle are handled directly by Paddle in accordance with their consumer protection rules and legal terms.
              </p>
              <p>
                By purchasing MiCartaPro, you agree to Paddle's Terms & Conditions and this policy.
              </p>
              <p className="mt-4">
                🔗 <strong>Paddle Terms:</strong> <a href="https://paddle.com/legal" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:text-blue-700 underline">https://paddle.com/legal</a>
              </p>
              <p>
                🔗 <strong>Paddle Privacy:</strong> <a href="https://paddle.com/privacy" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:text-blue-700 underline">https://paddle.com/privacy</a>
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">1. 14-Day Refund Policy</h2>
              <p>
                If you are not satisfied with your purchase, you may request a refund within 14 days of the original transaction.
              </p>
              <p>
                Refunds are processed directly by Paddle and returned to the original payment method.
              </p>
              <p className="mt-4">
                To request a refund you may contact:
              </p>
              <p>
                <strong>MiCartaPro Support:</strong> support@micartapro.com
              </p>
              <p>
                <strong>Paddle Support:</strong> <a href="https://paddle.net/support" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:text-blue-700 underline">https://paddle.net/support</a>
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">2. Subscription Renewals</h2>
              <p>
                Subscriptions renew automatically unless cancelled.
              </p>
              <p>
                You may cancel at any time prior to the next billing date to stop future charges.
              </p>
              <p>
                Cancelling a subscription prevents future payments — it does not automatically generate a refund for past charges.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">3. Human-Driven Add-on Services (Separate Billing)</h2>
              <p>
                MiCartaPro may offer optional services such as:
              </p>
              <ul className="list-disc pl-6 space-y-2 mb-4">
                <li>Menu setup / customization</li>
                <li>Data import / migration</li>
                <li>Consulting or manual implementation</li>
              </ul>
              <p>
                These add-on services are billed separately and are not processed through Paddle, therefore refunds for these services are handled directly between MiCartaPro and the client under separate terms.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">4. Tax & Billing</h2>
              <p>
                Paddle may collect and remit VAT/GST or other taxes depending on your region.
              </p>
              <p>
                Invoices are issued electronically at purchase.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">5. Contact & Support</h2>
              <p>
                📧 Email: support@micartapro.com
              </p>
              <p>
                📱 +56 9 5785 7558
              </p>
              <p>
                For Paddle billing or refund requests: <a href="https://paddle.net/support" target="_blank" rel="noopener noreferrer" className="text-blue-600 hover:text-blue-700 underline">https://paddle.net/support</a>
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

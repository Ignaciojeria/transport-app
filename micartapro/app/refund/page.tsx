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
            <strong>Last updated:</strong> {new Date().toLocaleDateString('en-US', { year: 'numeric', month: 'long', day: 'numeric' })}
          </p>

          <div className="prose prose-lg max-w-none space-y-8 text-gray-700">
            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">1. Overview</h2>
              <p>
                This Refund Policy outlines the terms and conditions for refunds of MiCartaPro services. Our payment processing is handled by Paddle, who acts as the Merchant of Record and authorized reseller of our Product. By purchasing our services, you agree to the terms outlined in this Policy and Paddle's Invoiced Consumer Terms and Conditions.
              </p>
              <p>
                For detailed information about Paddle's refund policy, please refer to Paddle's Invoiced Consumer Terms and Conditions available at paddle.com/legal.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">2. Consumer Right to Cancel</h2>
              <p>
                If you are a Consumer (purchasing for personal use) and unless the exception below applies, you have the right to cancel this Agreement and return the Product within <strong>14 days</strong> without giving any reason. The cancellation period will expire after 14 days from the day after completion of the Transaction.
              </p>
              <p>
                To meet the cancellation deadline, it is sufficient that you send us your communication concerning your exercise of the cancellation right before the expiration of the 14 day period.
              </p>
              <p>
                To cancel your order, you must inform Paddle or us of your decision. To ensure immediate processing, please contact us via WhatsApp at +56957857558 or email at support@micartapro.com, or contact Paddle directly through their support channels.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">3. Exception to the Right to Cancel</h2>
              <p>
                Your right as a Consumer to cancel your order <strong>does not apply</strong> to the supply of Digital Content that you have started to download, stream or otherwise acquire and to Products which you have had the benefit of.
              </p>
              <p>
                Since our digital menu service (MiCartaPro) is immediately made available upon purchase, by completing your purchase and accessing the service, you consent to immediate performance of this Agreement and acknowledge that you will lose your right of withdrawal from this Agreement once you have started to use the service.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">4. Effect of Cancellation</h2>
              <p>
                If you cancel this Agreement as permitted above, we will reimburse to you all payments received from you.
              </p>
              <p>
                We will make the reimbursement without undue delay, and not later than 14 days after the day on which we are informed about your decision to cancel this Agreement.
              </p>
              <p>
                We will make the reimbursement using the same means of payment as you used for the initial transaction and you will not incur any fees as a result of the reimbursement.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">5. Refund Policy</h2>
              <p>
                Refunds are provided at the sole discretion of Paddle and on a case-by-case basis and may be refused. Paddle will refuse a refund request if they find evidence of fraud, refund abuse, or other manipulative behaviour that entitles Paddle to counterclaim the refund.
              </p>
              <p>
                This does not affect your rights as a Consumer in relation to Products which are not as described, faulty or not fit for purpose.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">6. Subscriptions</h2>
              <p>
                Our services are provided on a subscription basis ("Paid Subscriptions"). Paid Subscriptions automatically renew until cancelled.
              </p>
              <p>
                <strong>Important:</strong> In respect of subscription services, your right to cancel is only present following the initial subscription and not upon each automatic renewal.
              </p>
              <p>
                If you wish to cancel your subscription, please contact us at least 48 hours before the end of the current billing period. Please make sure you provide your order number and the email address used to purchase the Product. Your cancellation will take effect at the next payment date.
              </p>
              <p>
                <strong>There are no refunds on unused subscription periods.</strong>
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">7. Payment Processing</h2>
              <p>
                All payments are processed by Paddle, who acts as the Merchant of Record. Paddle handles all payment transactions, refunds, and related customer service inquiries. For questions about payment processing, please refer to Paddle's terms and conditions at paddle.com/legal.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">8. Contact</h2>
              <p>
                If you have questions about this Refund Policy or need to request a refund, please contact us:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>WhatsApp:</strong> +56957857558</li>
                <li><strong>Email:</strong> support@micartapro.com</li>
                <li><strong>Paddle Support:</strong> You can also contact Paddle directly through their support channels</li>
              </ul>
              <p className="mt-4">
                We aim to respond to all refund requests and inquiries within 2-3 business days.
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

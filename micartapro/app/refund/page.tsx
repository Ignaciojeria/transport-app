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
                This Refund Policy ("Policy") outlines the terms and conditions under which MiCartaPro ("we", "our", or "the company") provides refunds for our digital menu services. By purchasing our services, you agree to the terms outlined in this Policy.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">2. Service Description</h2>
              <p>
                MiCartaPro provides digital menu services including:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Custom digital menu design and development</li>
                <li>Exclusive QR code generation</li>
                <li>Shopping cart integration</li>
                <li>WhatsApp integration for order reception</li>
                <li>Responsive design for all devices</li>
                <li>First year of service included in initial payment</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">3. Pricing Structure</h2>
              <p>
                Our services are priced as follows:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>Initial Setup:</strong> Starting from $150 USD (one-time payment)</li>
                <li><strong>First Year:</strong> Included in initial payment (free promotional period)</li>
                <li><strong>Renewal (Second Year Onwards):</strong> $10 USD per month</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">4. Refund Eligibility</h2>
              
              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">4.1. Initial Setup Fee</h3>
              <p>
                Refunds for the initial setup fee ($150 USD) may be requested under the following circumstances:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>Before Service Commencement:</strong> If you request a refund before we begin work on your digital menu, you may receive a full refund within 7 days of payment.</li>
                <li><strong>Service Not Delivered:</strong> If we fail to deliver the agreed-upon services within the specified timeframe and you have not approved any delays, you may be eligible for a full or partial refund.</li>
                <li><strong>Technical Issues:</strong> If the delivered service has critical technical issues that we cannot resolve within 30 days, you may be eligible for a refund.</li>
              </ul>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">4.2. Monthly Renewal Fees</h3>
              <p>
                Monthly renewal fees ($10 USD) are non-refundable once paid. However, you may cancel your subscription at any time to avoid future charges. Cancellation will take effect at the end of the current billing period.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">5. Non-Refundable Situations</h2>
              <p>
                Refunds will not be provided in the following situations:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>After the service has been fully delivered and accepted by you</li>
                <li>If you have used the service for more than 30 days after delivery</li>
                <li>If the service was modified or customized according to your specific requirements and you have approved the final version</li>
                <li>If you have violated our Terms and Conditions, resulting in service termination</li>
                <li>If you request a refund due to a change of mind after the service has been delivered</li>
                <li>For monthly renewal fees that have already been charged</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">6. Refund Process</h2>
              
              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">6.1. Requesting a Refund</h3>
              <p>
                To request a refund, you must:
              </p>
              <ol className="list-decimal pl-6 space-y-2">
                <li>Contact us via WhatsApp at +56957857558 or email at support@micartapro.com</li>
                <li>Provide your order details, including payment confirmation and service information</li>
                <li>Clearly state the reason for your refund request</li>
                <li>Include any relevant documentation or evidence supporting your request</li>
              </ol>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">6.2. Refund Review</h3>
              <p>
                We will review your refund request within 5-7 business days. During this time, we may:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Request additional information or clarification</li>
                <li>Investigate the circumstances of your request</li>
                <li>Attempt to resolve any issues that may have led to the refund request</li>
              </ul>

              <h3 className="text-xl font-semibold text-gray-900 mb-3 mt-4">6.3. Refund Approval and Processing</h3>
              <p>
                If your refund request is approved:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>We will process the refund within 10-15 business days</li>
                <li>The refund will be issued to the original payment method used for the purchase</li>
                <li>You will receive a confirmation email once the refund has been processed</li>
                <li>The refund amount will be based on the circumstances and may be partial or full</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">7. Partial Refunds</h2>
              <p>
                In some cases, we may offer a partial refund instead of a full refund. This may occur when:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Work has already been completed on your project</li>
                <li>You have used the service for a portion of the paid period</li>
                <li>Some deliverables have been accepted and used</li>
                <li>We determine that a partial refund is fair based on the services already rendered</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">8. Cancellation Policy</h2>
              <p>
                You may cancel your subscription at any time:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>Monthly Subscriptions:</strong> Cancel before the next billing cycle to avoid future charges. No refunds for the current period.</li>
                <li><strong>Service Cancellation:</strong> If you cancel during the initial setup phase, refund eligibility will be determined based on the work completed.</li>
                <li><strong>Immediate Cancellation:</strong> Upon cancellation, your service will remain active until the end of the current billing period.</li>
              </ul>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">9. Disputes and Chargebacks</h2>
              <p>
                If you have a dispute regarding a charge, we encourage you to contact us directly before initiating a chargeback with your payment provider. Chargebacks may result in:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li>Immediate suspension of your account</li>
                <li>Termination of services</li>
                <li>Additional fees or penalties</li>
              </ul>
              <p className="mt-4">
                We are committed to resolving disputes fairly and promptly through direct communication.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">10. Currency and Payment Method</h2>
              <p>
                All refunds will be processed in the same currency as the original payment (USD). Refunds will be issued to the original payment method used for the purchase. Processing times may vary depending on your payment provider.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">11. Changes to This Policy</h2>
              <p>
                We reserve the right to modify this Refund Policy at any time. Changes will be effective immediately upon posting on this page. We will notify you of any material changes via email or through our services.
              </p>
              <p>
                Your continued use of our services after changes to this Policy constitutes acceptance of the updated terms.
              </p>
            </section>

            <section>
              <h2 className="text-2xl font-semibold text-gray-900 mb-4">12. Contact</h2>
              <p>
                If you have questions about this Refund Policy or need to request a refund, please contact us:
              </p>
              <ul className="list-disc pl-6 space-y-2">
                <li><strong>WhatsApp:</strong> +56957857558</li>
                <li><strong>Email:</strong> support@micartapro.com</li>
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

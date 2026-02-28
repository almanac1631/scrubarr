document.addEventListener('DOMContentLoaded', () => {
    const tooltip = document.getElementById('global-tooltip');

    document.addEventListener('mouseover', (e) => {
        const trigger = e.target.closest('[data-tooltip]');

        if (!trigger) {
            tooltip.classList.add('hidden', 'opacity-0');
            return;
        }

        const tooltipKey = trigger.dataset.tooltip;

        if (tooltipKey === "status-info") {
            tooltip.innerHTML = "";
            const decision = getDecisionStr(trigger.dataset.decision);
            const torrentStatus = trigger.dataset.torrentStatus;
            const torrentRatio = trigger.dataset.torrentRatio;
            const torrentAge = trigger.dataset.torrentAge;

            const trackerName = trigger.dataset.trackerName;
            const trackerMinRatio = formatFloatStr(trigger.dataset.trackerMinRatio);
            const trackerMinAge = trigger.dataset.trackerMinAge;

            const decisionElem = document.createElement("div");
            decisionElem.classList.add("font-bold");
            decisionElem.textContent = decision;
            tooltip.append(decisionElem);

            if (torrentStatus === "present") {
                if (trackerName !== "") {
                    const trackerNameElem = document.createElement("div");
                    trackerNameElem.textContent = trackerName;
                    tooltip.append(trackerNameElem);
                }
                if (torrentRatio !== "-1") {
                    const ratioElem = document.createElement("div");
                    ratioElem.textContent = `Ratio: ${formatFloatStr(torrentRatio)}`;
                    if (trackerMinRatio !== "") {
                        ratioElem.textContent += `/${formatFloatStr(trackerMinRatio)}`;
                    }
                    tooltip.append(ratioElem);
                }
                if (torrentAge !== "-1") {
                    const ageElem = document.createElement("div");
                    ageElem.textContent = `Age: ${nanosecondsToDays(Number(torrentAge))}d`;
                    if (trackerMinAge !== "") {
                        ageElem.textContent += `/${nanosecondsToDays(Number(trackerMinAge))}d`;
                    }
                    tooltip.append(ageElem);
                }
            } else {
                const torrentInfoElem = document.createElement("div");
                torrentInfoElem.textContent = "No torrent entry";
                tooltip.append(torrentInfoElem);
            }
        } else if (tooltipKey === "disk-quota") {
            const diskQuotaUsed = trigger.dataset.diskQuotaUsed;
            const diskQuotaFree = trigger.dataset.diskQuotaFree;
            tooltip.innerHTML = "";

            const usedElem = document.createElement("div");
            usedElem.textContent = `Used: ${diskQuotaUsed}`;
            tooltip.append(usedElem);

            const freeElem = document.createElement("div");
            freeElem.textContent = `Free: ${diskQuotaFree}`;
            tooltip.append(freeElem);
        } else {
            return;
        }

        const rect = trigger.getBoundingClientRect();
        tooltip.classList.remove('hidden');
        let top = rect.top - tooltip.offsetHeight - 8; // 8px gap
        if (trigger.dataset.tooltipPosition === "bottom") {
            top = rect.bottom + 8;
        }
        const left = rect.left + (rect.width / 2) - (tooltip.offsetWidth / 2);
        tooltip.style.top = `${top}px`;
        tooltip.style.left = `${left}px`;

        requestAnimationFrame(() => {
            tooltip.classList.remove('opacity-0');
        });
    });

    document.addEventListener('mouseout', (e) => {
        const trigger = e.target.closest('[data-tooltip]');
        if (trigger) {
            tooltip.classList.add('opacity-0');
        }
    });
});

const nanosecondsToDays = (ns) => Math.floor(ns / 86_400_000_000_000);

const formatFloatStr = (floatStr) => new Intl.NumberFormat('en-US', {
    minimumFractionDigits: 0,
    maximumFractionDigits: 2
}).format(floatStr);

function getDecisionStr(decision) {
    if (decision === "safe_to_delete") {
        return "Safe to delete"
    } else if (decision === "pending") {
        return "Status pending"
    } else {
        return "???"
    }
}